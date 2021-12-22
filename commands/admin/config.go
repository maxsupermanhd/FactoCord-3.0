package admin

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

var ConfigCommandDoc = support.Command{
	Name:    "config",
	Desc:    "Manage config.json",
	Usage:   "/config save | load | get <path> | set <path> <value>?",
	Doc:     "command manages FactoCord's config",
	Admin:   true,
	Command: nil, // ConfigCommand
	Subcommands: []support.Command{
		{
			Name:  "save",
			Admin: true,
			Desc:  "Save FactoCord's config from memory to `config.json`",
			Doc: "command saves FactoCord's config from memory to `config.json`.\n" +
				"It also adds any options missing from config.json",
			Command: respond(save),
		},
		{
			Name:  "load",
			Admin: true,
			Desc:  "Load the config from `config.json`",
			Doc: "command loads the config from `config.json`.\n" +
				"Any unsaved changes after the last `/config save` command will be lost.",
			Command: respond(load),
		},
		{
			Name:  "get",
			Admin: true,
			Usage: "/config get <path>?",
			Desc:  "Get the value of a config setting specified by `path`",
			Doc: "command outputs the value of a config setting specified by <path>.\n" +
				"All path members are separated by a dot '.'\n" +
				"If the path is empty, it outputs the whole config.\n" +
				"Secrets like discord_token are kept secret.\n" +
				"Examples:\n" +
				"```\n" +
				"/config get\n" +
				"/config get admin_ids\n" +
				"/config get admin_ids.0\n" +
				"/config get command_roles\n" +
				"/config get command_roles.mod\n" +
				"/config get messages\n" +
				"```",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "path",
					Description: "Path of the setting",
					Required:    false,
				},
			},
			Command: respond(get),
		},
		{
			Name:  "set",
			Admin: true,
			Usage: "/config set <path>\n" +
				"/config set <path> <value>",
			Desc: "Set or delete the value of a config setting specified by `path`",
			Doc: "command sets the value of a config setting specified by <path>.\n" +
				"This command can set only simple types such as strings, numbers, and booleans.\n" +
				"If no value is specified, this command deletes the value if possible, otherwise it sets it to a zero-value (0, \"\", false).\n" +
				"To add a value to an array or an object specify it's index as `*` (e.g. `/config set admin_ids.* 1234`).\n" +
				"Changes made by this command are not automatically saved. Use `/config save` to do it.\n" +
				"Examples:" +
				"```\n" +
				"/config set game_name \"Factorio 1.0\"\n" +
				"/config set ingame_discord_user_colors true\n" +
				"/config set admin_ids.0 123456789\n" +
				"/config set admin_ids.* 987654321\n" +
				"/config set command_roles.mod 55555555\n" +
				"/config set messages.server_save **:mango: Game saved!**\n" +
				"```",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "path",
					Description: "Path of the setting",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "value",
					Description: "Value to set. If unspecified, deletes the setting.",
					Required:    false,
				},
			},
			Command: respond(set),
		},
	},
}

func respond(f func(i *discordgo.InteractionCreate) string) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		support.Respond(s, i, f(i))
	}
}

func save(_ *discordgo.InteractionCreate) string {
	str, err := json.MarshalIndent(support.Config, "", "    ")
	if err != nil {
		support.Panik(err, "... when converting config to json")
		return "Error when converting config to json"
	}
	err = ioutil.WriteFile(support.ConfigPath, str, 0666)
	if err != nil {
		support.Panik(err, "... when saving config.json")
		return "Error when saving config.json"
	}
	return "Config saved"
}

func load(_ *discordgo.InteractionCreate) string {
	err := support.Config.Load()
	if err != nil {
		return err.Error()
	}
	return "Config reloaded"
}

func get(i *discordgo.InteractionCreate) string {
	path := ""
	options := i.ApplicationCommandData().Options[0].Options
	if len(options) == 1 {
		path = options[0].StringValue()
	}
	if strings.Contains(path, " \n\t") {
		return "There should be no spaces in the path"
	}
	var value interface{}
	if path == "" {
		config := support.Config // copy
		config.DiscordToken = "my precious"
		config.Username = "my precious"
		config.ModPortalToken = "my precious"
		value = config
	} else {
		path := strings.Split(path, ".")
		if path[0] == "discord_token" {
			return "Shhhh, it's a secret"
		}
		x, err := walk(&support.Config, path)
		if err != nil {
			return err.Error()
		}
		value = x.Interface()
	}
	res, err := json.MarshalIndent(value, "", "    ")
	if err != nil {
		support.Panik(err, "... when converting to json")
		return "Error when converting to json"
	}
	return fmt.Sprintf("```json\n%s\n```", string(res))
}

func set(i *discordgo.InteractionCreate) string {
	options := i.ApplicationCommandData().Options[0].Options
	pathS := options[0].StringValue()
	valueS := ""
	if len(options) == 2 {
		valueS = options[1].StringValue()
	}
	if pathS == "" {
		return "Path is empty"
	}
	path := strings.Split(pathS, ".")
	if path[0] == "discord_token" {
		return "Are trying to brainwash me?"
	}
	name := path[len(path)-1]
	pathTo := strings.Join(path[:len(path)-1], ".")
	if pathTo == "" {
		pathTo = "."
	}
	current, err := walk(&support.Config, path[:len(path)-1])
	if err != nil {
		return err.Error()
	}
	switch current.Kind() {
	case reflect.Slice:
		if name == "*" {
			value, errs := createValue(current.Type().Elem(), valueS)
			if errs != "" {
				return pathS + errs
			}
			current.Set(reflect.Append(current, value))
		} else {
			num, err := strconv.ParseUint(name, 10, 0)
			if err != nil {
				return fmt.Sprintf("%s is array but \"%s\" is not an int", pathS, name)
			}
			if current.Len() <= int(num) {
				return fmt.Sprintf("%d is bigger than %s's size (%d)", num, pathS, current.Len())
			}
			if valueS == "" {
				sliceRemove(current, int(num))
			} else {
				value, errs := createValue(current.Type().Elem(), valueS)
				if errs != "" {
					return pathS + errs
				}
				current.Index(int(num)).Set(value)
			}
		}
	case reflect.Struct:
		fieldName := getFieldByTag(name, "json", current.Type())
		if fieldName == "" {
			if pathTo == "." {
				return fmt.Sprintf("config does not have an option called \"%s\"", name)
			} else {
				return fmt.Sprintf("struct %s does not have a field called \"%s\"", pathTo, name)
			}
		}
		field := current.FieldByName(fieldName)
		value, errs := createValue(field.Type(), valueS)
		if errs != "" {
			return pathS + errs
		}
		field.Set(value)
	case reflect.Map:
		key, errs := createValue(current.Type().Key(), name)
		if errs != "" {
			return pathS + errs
		}
		var value reflect.Value
		if valueS == "" {
			value = reflect.Value{}
		} else {
			value, errs = createValue(current.Type().Elem(), valueS)
			if errs != "" {
				return pathS + errs
			}
		}
		current.SetMapIndex(key, value)
	default:
		return fmt.Sprintf("%s's type (%s) is not supported", pathS, current.Type().String())
	}
	return "Value set"
}

func walk(v interface{}, path []string) (reflect.Value, error) {
	var current = reflect.ValueOf(v)
	if current.Type().Kind() != reflect.Ptr {
		panic("walk: v should be pointer")
	}
	current = current.Elem()
	for i, name := range path {
		walkedPath := strings.Join(path[:i], ".")
		switch current.Kind() {
		case reflect.Slice:
			num, err := strconv.ParseUint(name, 10, 0)
			if err != nil {
				return reflect.Value{}, fmt.Errorf("%s is array but \"%s\" is not an int", walkedPath, name)
			}
			if current.Len() <= int(num) {
				return reflect.Value{}, fmt.Errorf("size of %s[%d] is less than %d", walkedPath, current.Len(), num)
			}
			current = current.Index(int(num))
		case reflect.Struct:
			field := getFieldByTag(name, "json", current.Type())
			if field == "" {
				return reflect.Value{}, fmt.Errorf("struct %s does not have a field called \"%s\"", walkedPath, name)
			}
			current = current.FieldByName(field)
		case reflect.Map:
			key, errs := createValue(current.Type().Key(), name)
			if errs != "" {
				return reflect.Value{}, fmt.Errorf(walkedPath + "." + name + errs)
			}
			current = current.MapIndex(key)
			if !current.IsValid() {
				return reflect.Value{}, fmt.Errorf("%s does not have key \"%s\"", walkedPath, name)
			}
		default:
			return reflect.Value{}, fmt.Errorf("%s's type (%s) is not supported", walkedPath, current.Type().String())
		}
	}
	return current, nil
}

func createValue(t reflect.Type, value string) (reflect.Value, string) {
	if value == "" {
		return reflect.New(t).Elem(), ""
	}
	switch t.Kind() {
	case reflect.Bool:
		val, err := strconv.ParseBool(value)
		if err != nil {
			return reflect.Value{}, fmt.Sprintf(" requires bool but \"%s\" can't be converted to bool", value)
		}
		return reflect.ValueOf(val), ""
	case reflect.Int:
		num, err := strconv.ParseUint(value, 10, 0)
		if err != nil {
			return reflect.Value{}, fmt.Sprintf(" requires int but \"%s\" is not an int", value)
		}
		return reflect.ValueOf(int(num)), ""
	case reflect.String:
		if value[0] == '"' && value[len(value)-1] == '"' {
			return reflect.ValueOf(value[1 : len(value)-1]), ""
		}
		return reflect.ValueOf(value), ""
	default:
		return reflect.Value{}, fmt.Sprintf("'s type (%s) is not supported", t.String())
	}
}

func getFieldByTag(tag, key string, t reflect.Type) (fieldname string) {
	if t.Kind() != reflect.Struct {
		panic("bad type")
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		v := strings.Split(f.Tag.Get(key), ",")[0] // use split to ignore tag "options" like omitempty, etc.
		if v == tag {
			return f.Name
		}
	}
	return ""
}

func sliceRemove(v reflect.Value, index int) {
	for i := index; i < v.Len()-1; i++ {
		v.Index(i).Set(v.Index(i + 1))
	}
	v.SetLen(v.Len() - 1)
}
