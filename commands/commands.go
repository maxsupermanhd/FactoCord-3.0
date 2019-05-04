package commands

import (
	"strings"

	"./admin"
	"./utils"
	"../support"
	"github.com/bwmarrin/discordgo"
)

// Commands is a struct containing a slice of Command.
type Commands struct {
	CommandList []Command
}

// Command is a struct containing fields that hold command information.
type Command struct {
	Name    string
	Command func(s *discordgo.Session, m *discordgo.MessageCreate)
        Command1 func(s *discordgo.Session, m *discordgo.MessageCreate, arg1 string)
        Command2 func(s *discordgo.Session, m *discordgo.MessageCreate, arg1 string, arg2 string)
	Admin   bool
	Args	int
	Desc	string
}

// CL is a Commands interface.
var CL Commands

// RegisterCommands registers the commands on start up.
func RegisterCommands() {
	// Admin Commands
	CL.CommandList = append(CL.CommandList, Command{Name: "Stop", Command: admin.StopServer,
		Admin: true, Args:0, Desc: "Save the game and stop the factorio server."})
	CL.CommandList = append(CL.CommandList, Command{Name: "Restart", Command: admin.Restart,
		Admin: true, Args:0, Desc: "Save the game and restart the factorio server."})
        CL.CommandList = append(CL.CommandList, Command{Name: "Save", Command: admin.SaveServer,
                Admin: true, Args:0, Desc: "Save the game."})
	CL.CommandList = append(CL.CommandList, Command{Name: "Kick", Command2: admin.KickPlayer,
		Admin: true, Args:2, Desc: "Kick a user from the server. Usage: !kick <player> <reason>"})
	CL.CommandList = append(CL.CommandList, Command{Name: "Ban", Command2: admin.BanPlayer,
		Admin: true, Args:2, Desc: "Ban a user from the server. Usage: !ban <player> <reason>"})
	CL.CommandList = append(CL.CommandList, Command{Name: "Unban", Command1: admin.UnbanPlayer,
		Admin: true, Args:1, Desc: "Unban a user from the server. Usage !unban <player>"})
	// Util Commands
	CL.CommandList = append(CL.CommandList, Command{Name: "Mods", Command: utils.ModsList,
		Admin: false, Args:0, Desc: "List the mods on the server"})
        CL.CommandList = append(CL.CommandList, Command{Name: "List",
                Admin: false, Args:0, Desc: "List the commands for Factocord"})
}

func commandListEmbed() *discordgo.MessageEmbed {
	fields := []*discordgo.MessageEmbedField{}
	for _, command := range CL.CommandList {
		strAdmin := ""
		if command.Admin == true {
			strAdmin = " - Admin Only!"
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name: "!"+ command.Name,
			Value: command.Desc + strAdmin,
		})
	}
	embed := &discordgo.MessageEmbed{
		Type:  "rich",
		Color: 52,
		Description: "List of all commands currently available in version 2.0 of FactoCord",
		Title:  "FactoCord Commands",
		Fields: fields,
	}
	return embed
}

// RunCommand runs a specified command.
func RunCommand(input string, s *discordgo.Session, m *discordgo.MessageCreate) {
	inputvars := strings.SplitN(input, " ", 3)
	for _, command := range CL.CommandList {
		if strings.ToLower(command.Name) == strings.ToLower("List") {
			s.ChannelMessageSendEmbed(support.Config.FactorioChannelID, commandListEmbed())
			return
		}

		if strings.ToLower(command.Name) == strings.ToLower(inputvars[0]) && command.Args == 0 && len(inputvars) == 1{
			//s.ChannelMessageSend(support.Config.FactorioChannelID, "0 arguments!")
			if command.Admin && CheckAdmin(m.Author.ID) {
				command.Command(s, m)
			}

			if !command.Admin {
				command.Command(s, m)
			}
			return
		} else if strings.ToLower(command.Name) == strings.ToLower(inputvars[0]) && command.Args == 1 && len(inputvars) == 2 {
                        //s.ChannelMessageSend(support.Config.FactorioChannelID, "1 arguments!")
			if command.Admin && CheckAdmin(m.Author.ID) {
                                command.Command1(s, m, inputvars[1])
                        }

                        if !command.Admin {
                                command.Command1(s, m, inputvars[1])
                        }
                        return
                }else if strings.ToLower(command.Name) == strings.ToLower(inputvars[0]) && command.Args == 2 && len(inputvars) == 3 {
			//s.ChannelMessageSend(support.Config.FactorioChannelID, "2 arguments!")
			if command.Admin && CheckAdmin(m.Author.ID) {
                                command.Command2(s, m, inputvars[1], inputvars[2])
                        }

                        if !command.Admin {
                                command.Command2(s, m, inputvars[1], inputvars[2])
                        }
			return
		}
	}
	s.ChannelMessageSend(support.Config.FactorioChannelID, "Command not found or the command is missing required arguments!")
}

// CheckAdmin checks if the user attempting to run an admin command is an admin
func CheckAdmin(ID string) bool {
	for _, admin := range support.Config.AdminIDs {
		if ID == admin {
			return true
		}
	}
	return false
}
