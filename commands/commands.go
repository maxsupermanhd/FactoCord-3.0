package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/maxsupermanhd/FactoCord-3.0/commands/admin"
	"github.com/maxsupermanhd/FactoCord-3.0/commands/utils"
	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

// Commands is a list of all available commands
var commandsList = [...]support.Command{
	// Admin Commands
	admin.ServerCommandDoc,
	admin.SaveServerDoc,
	admin.KickPlayerDoc,
	admin.BanPlayerDoc,
	admin.UnbanPlayerDoc,
	admin.ConfigCommandDoc,
	admin.ModCommandDoc,

	// Util Commands
	utils.ModListDoc,
	utils.VersionDoc,
	utils.InfoDoc,
	utils.OnlineDoc,
	{
		Name: "help",
		Desc: "List the commands for FactoСord and get documentation on commands and subcommands. Try `/help help`",
		Usage: "/help\n" +
			"/help <command>\n" +
			"/help <command> <subcommand>",
		Doc:     "command returns list of all commands and documentation about any command and its' subcommands",
		Command: nil, // = helpCommand in the init()
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "command",
				Description: "The command to get the help on",
				Required:    false,
				Choices:     nil,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "subcommand",
				Description: "The subcommand to get the help on",
				Required:    false,
			},
		},
	},
}
var Commands = map[string]support.Command{}

func init() {
	help := &commandsList[len(commandsList)-1]
	help.Command = helpCommand
	for _, command := range commandsList {
		Commands[command.Name] = command
		help.Options[0].Choices = append(help.Options[0].Choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  command.Name,
			Value: command.Name,
		})
	}
}

func helpCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if len(i.ApplicationCommandData().Options) == 0 {
		helpAllCommands(s, i)
		return
	}
	commandName := strings.ToLower(i.ApplicationCommandData().Options[0].StringValue())
	subcommand := ""
	if len(i.ApplicationCommandData().Options) == 2 {
		subcommand = strings.ToLower(i.ApplicationCommandData().Options[1].StringValue())
	}
	for _, command := range Commands {
		if command.Name == commandName {
			helpOnCommand(s, i, &command, subcommand)
			return
		}
	}
	support.Send(s, "There's no such command as \""+commandName+"\"")
}

func helpOnCommand(s *discordgo.Session, i *discordgo.InteractionCreate, command *support.Command, subcommandName string) {
	path := "/" + command.Name
	if subcommandName != "" {
		found := false
		for _, subcommand := range command.Subcommands {
			if subcommand.Name == subcommandName {
				command = &subcommand
				path += " " + subcommandName
				found = true
				break
			}
		}
		if !found {
			support.Send(s, fmt.Sprintf(`Command "%s" has no subcommand "%s"`, command.Name, subcommandName))
			return
		}
	}
	quoted := "`" + path + "`"
	embed := &discordgo.MessageEmbed{
		Type:        "rich",
		Color:       0x6289FF,
		Title:       fmt.Sprintf("Documentation on `%s` command", path),
		Description: quoted + " " + command.Doc,
	}
	usage := command.Usage
	if usage == "" {
		usage = path
	}
	if strings.Contains(usage, "\n") {
		usage = "```\n" + usage + "\n```"
	} else {
		usage = "`" + usage + "`"
	}
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:  "**Usage:**",
		Value: usage,
	})
	if len(command.Subcommands) > 0 {
		subcommands := ""
		for _, subcommand := range command.Subcommands {
			if subcommands != "" {
				subcommands += "\n"
			}
			subcommands += path + " " + subcommand.Name
		}
		subcommands = "```\n" + subcommands + "\n```"
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "**Subcommands:**",
			Value: subcommands,
		})
	}
	support.RespondComplex(s, i, &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{embed},
		Flags:  1 << 6, // EPHEMERAL
	})
}

func helpAllCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	fields := make([]*discordgo.MessageEmbedField, 0, len(Commands))

	for _, command := range commandsList {
		desc := command.Desc
		if roleID, exists := support.Config.CommandRoles[strings.ToLower(command.Name)]; exists {
			roles, err := s.GuildRoles(support.GuildID)
			if err != nil {
				support.Panik(err, "... when querying guild roles")
				return
			}
			found := false
			for _, role := range roles {
				if role.ID == roleID {
					found = true
					desc = "[Role \"" + role.Name + "\"] " + desc
					break
				}
			}
			if !found {
				desc = "[Role not found in guild] " + desc
			}
		} else if command.Admin {
			desc = "[Admin] " + desc
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "/" + command.Name,
			Value: desc,
		})
	}
	embed := &discordgo.MessageEmbed{
		Type:        "rich",
		Color:       52,
		Description: "List of all commands currently available in this version of FactoCord",
		Title:       "FactoCord Commands",
		Fields:      fields,
	}
	support.RespondComplex(s, i, &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{embed},
		Flags:  1 << 6, // EPHEMERAL
	})
}

// ProcessInteraction runs a specified command.
func ProcessInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name
	command, ok := Commands[name]
	if !ok {
		//support.Send(s, "Command not found. Try using \"/help\"")
		return
	}
	subcommand := command
	if len(i.ApplicationCommandData().Options) > 0 &&
		i.ApplicationCommandData().Options[0].Type == discordgo.ApplicationCommandOptionSubCommand {
		found := command.Subcommand(i.ApplicationCommandData().Options[0].Name)
		if found != nil {
			subcommand = *found
		}
	}

	execute := false
	err := ""

	if subcommand.Admin {
		if CheckAdmin(i.Member.User.ID) {
			execute = true
		} else {
			err = "You are not an admin!"
		}
	} else {
		execute = true
	}
	if roleID, exists := support.Config.CommandRoles[name]; exists {
		// TODO? role name
		err = "You don't have the required role"
		for _, memberRoleID := range i.Member.Roles {
			if memberRoleID == roleID {
				execute = true
			}
		}
	}
	if execute {
		if command.Command != nil {
			command.Command(s, i)
		} else {
			subcommand.Command(s, i)
		}
	} else {
		support.Respond(s, i, err)
	}
}

// CheckAdmin checks if the user attempting to run an admin command is an admin
func CheckAdmin(ID string) bool {
	// not using map here because they can be modified by /config set
	for _, adminID := range support.Config.AdminIDs {
		if ID == adminID {
			return true
		}
	}
	return false
}
