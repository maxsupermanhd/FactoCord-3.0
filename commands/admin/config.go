package admin

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/maxsupermanhd/FactoCord-3.0/v3/support"
)

var ConfigCommandDoc = support.CommandDoc{
	Name:  "config",
	Usage: "$config save | load | get <path> | set <path> <value>?",
	Doc:   "command manages FactoCord's config",
	Subcommands: []support.CommandDoc{
		{
			Name: "save",
			Doc: "command saves FactoCord's config from memory to `config.json`.\n" +
				"It also adds any options missing from config.json",
		},
		{
			Name: "load",
			Doc: "command loads the config from `config.json`.\n" +
				"Any unsaved changes after the last `$config save` command will be lost.",
		},
		{
			Name:  "get",
			Usage: "$config get <path>?",
			Doc: "command outputs the value of a config setting specified by <path>.\n" +
				"All path members are separated by a dot '.'\n" +
				"If the path is empty, it outputs the whole config.\n" +
				"Secrets like discord_token are kept secret.\n" +
				"Examples:\n" +
				"```\n" +
				"$config get\n" +
				"$config get admin_ids\n" +
				"$config get admin_ids.0\n" +
				"$config get command_roles\n" +
				"$config get command_roles.mod\n" +
				"$config get messages\n" +
				"```",
		},
		{
			Name: "set",
			Usage: "$config set <path>\n" +
				"$config set <path> <value>",
			Doc: "command sets the value of a config setting specified by <path>.\n" +
				"This command can set only simple types such as strings, numbers, and booleans.\n" +
				"If no value is specified, this command deletes the value if possible, otherwise it sets it to a zero-value (0, \"\", false).\n" +
				"To add a value to an array or an object specify it's index as '*' (e.g. `$config set admin_ids.* 1234`).\n" +
				"Changes made by this command are not automatically saved. Use `$config save` to do it.\n" +
				"Examples:" +
				"```\n" +
				"$config set prefix !\n" +
				"$config set game_name \"Factorio 1.0\"\n" +
				"$config set ingame_discord_user_colors true\n" +
				"$config set admin_ids.0 123456789\n" +
				"$config set admin_ids.* 987654321\n" +
				"$config set command_roles.mod 55555555\n" +
				"$config set messages.server_save **:mango: Game saved!**\n" +
				"```",
		},
	},
}

// ModCommand returns the list of mods running on the server.
func ConfigCommand(s *discordgo.Session, args string) {
	if args == "" {
		support.SendFormat(s, "Usage: "+ConfigCommandDoc.Usage)
		return
	}
	action, args := support.SplitDivide(args, " ")
	args = strings.TrimSpace(args)
	if _, ok := commands[action]; !ok {
		support.SendFormat(s, "Usage: "+ConfigCommandDoc.Usage)
		return
	}
	res := commands[action](args)
	support.Send(s, res)
}

var commands = map[string]func(string) string{
	"save": save,
	"load": load,
	"get":  get,
	"set":  set,
}

func save(args string) string {
	if args != "" {
		return "Save accepts no arguments"
	}
	s, err := json.MarshalIndent(support.Config, "", "    ")
	if err != nil {
		support.Panik(err, "... when converting config to json")
		return "Error when converting config to json"
	}
	err = os.WriteFile(support.ConfigPath, s, 0666)
	if err != nil {
		support.Panik(err, "... when saving config.json")
		return "Error when saving config.json"
	}
	return "Config saved"
}

func load(args string) string {
	if args != "" {
		return "Load accepts no arguments"
	}
	err := support.Config.Load()
	if err != nil {
		return err.Error()
	}
	return "Config reloaded"
}

func get(args string) string {
	if strings.Contains(args, " \n\t") {
		return "Why are there spaces in the path?"
	}
	var value interface{}
	if args == "" {
		config := support.Config // copy
		config.DiscordToken = "my precious"
		config.Username = "my precious"
		config.ModPortalToken = "my precious"
		value = config
	} else {
		path := strings.Split(args, ".")
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

func set(args string) string {
	pathS, valueS := support.SplitDivide(args, " ")
	if pathS == "" {
		return support.FormatUsage("Usage: $config set <path> <value>?")
	}
	path := strings.Split(pathS, ".")
	if len(path) == 0 {
		return "wtf??"
	}
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
