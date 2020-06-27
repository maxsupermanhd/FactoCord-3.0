package commands

import (
	"strings"

	"../support"
	"./admin"
	"./utils"
	"github.com/bwmarrin/discordgo"
)

// Command is a struct containing fields that hold command information.
type Command struct {
	Name string

	Command func(s *discordgo.Session, args string)

	Admin bool
	Desc  string
}

// Commands is a list of all available commands
var Commands = [...]Command{
	// Admin Commands
	{
		Name:    "Stop",
		Command: admin.StopServer,
		Admin:   true,
		Desc:    "Save the game and stop the factorio server.",
	},
	{
		Name:    "Restart",
		Command: admin.Restart,
		Admin:   true,
		Desc:    "Save the game and restart the factorio server.",
	},
	{
		Name:    "Save",
		Command: admin.SaveServer,
		Admin:   true,
		Desc:    "Save the game.",
	},
	{
		Name:    "Kick",
		Command: admin.KickPlayer,
		Admin:   true,
		Desc:    "Kick a user from the server. " + admin.KickPlayerUsage,
	},
	{
		Name:    "Ban",
		Command: admin.BanPlayer,
		Admin:   true,
		Desc:    "Ban a user from the server. " + admin.BanPlayerUsage,
	},
	{
		Name:    "Unban",
		Command: admin.UnbanPlayer,
		Admin:   true,
		Desc:    "Unban a user from the server. " + admin.UnbanPlayerUsage,
	},
	{
		Name:    "Config",
		Command: admin.ConfigCommand,
		Admin:   true,
		Desc:    "Manage config.json. " + admin.ConfigCommandUsage,
	},
	{
		Name:    "Mod",
		Command: admin.ModCommand,
		Admin:   true,
		Desc:    "Manage mod-list.json. " + admin.ModCommandUsage,
	},

	// Util Commands
	{
		Name:    "Mods",
		Command: utils.ModsList,
		Admin:   false,
		Desc:    "List the mods on the server. " + utils.ModListUsage,
	},
	{
		Name:    "Version",
		Command: utils.VersionString,
		Admin:   false,
		Desc:    "Get server version " + utils.VersionStringUsage,
	},
	{
		Name:  "Help",
		Admin: false,
		Desc:  "List the commands for Factocord",
	},
}

func helpCommand(s *discordgo.Session, m *discordgo.Message) {
	var fields []*discordgo.MessageEmbedField
	for _, command := range Commands {
		desc := support.FormatUsage(command.Desc)
		if roleID, exists := support.Config.CommandRoles[strings.ToLower(command.Name)]; exists {
			roles, err := s.GuildRoles(m.GuildID)
			if err != nil {
				support.Panik(err, "... when querying guild roles")
				return
			}
			found := false
			for _, role := range roles {
				if role.ID == roleID {
					found = true
					desc += " - Role \"" + role.Name + "\""
					break
				}
			}
			if !found {
				desc += " - Role not found in guild"
			}
		} else if command.Admin {
			desc += " - Admin Only!"
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  support.Config.Prefix + command.Name,
			Value: desc,
		})
	}
	embed := &discordgo.MessageEmbed{
		Type:        "rich",
		Color:       52,
		Description: "List of all commands currently available in version 3.0 of FactoCord",
		Title:       "FactoCord Commands",
		Fields:      fields,
	}
	support.SendEmbed(s, embed)
}

// RunCommand runs a specified command.
func RunCommand(input string, s *discordgo.Session, m *discordgo.Message) {
	inputvars := strings.SplitN(input+" ", " ", 2)
	commandName := strings.ToLower(inputvars[0])
	args := strings.TrimSpace(inputvars[1])

	if commandName == strings.ToLower("Help") {
		helpCommand(s, m)
		return
	}

	for _, command := range Commands {
		if strings.ToLower(command.Name) == commandName {
			execute := false
			err := ""

			if command.Admin {
				if CheckAdmin(m.Author.ID) {
					execute = true
				} else {
					err = "You are not an admin!"
				}
			} else {
				execute = true
			}
			if roleID, exists := support.Config.CommandRoles[commandName]; exists {
				// TODO? role name
				err = "You don't have the required role"
				for _, memberRoleID := range m.Member.Roles {
					if memberRoleID == roleID {
						execute = true
					}
				}
			}
			if execute {
				command.Command(s, args)
			} else {
				support.Send(s, err)
			}
			return
		}
	}
	support.SendFormat(s, "Command not found. Try using \"$help\"")
}

// CheckAdmin checks if the user attempting to run an admin command is an admin
func CheckAdmin(ID string) bool {
	for _, adminID := range support.Config.AdminIDs {
		if ID == adminID {
			return true
		}
	}
	return false
}
