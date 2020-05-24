package commands

import (
    "strings"

    "./admin"
    "./utils"
    "../support"
    "github.com/bwmarrin/discordgo"
)

// Command is a struct containing fields that hold command information.
type Command struct {
    Name    string

    Command func(s *discordgo.Session, m *discordgo.MessageCreate, args string)
    
    Admin   bool
    Desc	string
}


var Commands = [...]Command{
    // Admin Commands
    {
        Name: "Stop",
        Command: admin.StopServer,
        Admin: true,
        Desc: "Save the game and stop the factorio server.",
    },
    {
        Name: "Restart",
        Command: admin.Restart,
        Admin: true,
        Desc: "Save the game and restart the factorio server.",
    },
    {
        Name: "Save",
        Command: admin.SaveServer,
        Admin: true,
        Desc: "Save the game.",
    },
    {
        Name: "Kick",
        Command: admin.KickPlayer,
        Admin: true,
        Desc: "Kick a user from the server. " + admin.KickPlayerUsage,
    },
    {
        Name: "Ban",
        Command: admin.BanPlayer,
        Admin: true,
        Desc: "Ban a user from the server. " + admin.BanPlayerUsage,
    },
    {
        Name: "Unban",
        Command: admin.UnbanPlayer,
        Admin: true,
        Desc: "Unban a user from the server. " + admin.UnbanPlayerUsage,
    },
    
    // Util Commands
    {
        Name: "Mods",
        Command: utils.ModsList,
        Admin: false,
        Desc: "List the mods on the server. " + utils.ModListUsage,
    },
    {
        Name: "Help",
        Admin: false,
        Desc: "List the commands for Factocord",
    },
}

func commandListEmbed() *discordgo.MessageEmbed {
    fields := []*discordgo.MessageEmbedField{}
    for _, command := range Commands {
        desc := support.FormatUsage(command.Desc)
        if command.Admin {
            desc += " - Admin Only!"
        }
        fields = append(fields, &discordgo.MessageEmbedField{
            Name: support.Config.Prefix + command.Name,
            Value: desc,
        })
    }
    embed := &discordgo.MessageEmbed{
        Type:  "rich",
        Color: 52,
        Description: "List of all commands currently available in version 3.0 of FactoCord",
        Title:  "FactoCord Commands",
        Fields: fields,
    }
    return embed
}

// RunCommand runs a specified command.
func RunCommand(input string, s *discordgo.Session, m *discordgo.MessageCreate) {
    inputvars := strings.SplitN(input + " ", " ", 2)
    command_name := inputvars[0]
    args := strings.TrimSpace(inputvars[1])

    for _, command := range Commands {
        if strings.ToLower(command.Name) == strings.ToLower("Help") {
            s.ChannelMessageSendEmbed(support.Config.FactorioChannelID, commandListEmbed())
            return
        }

        if strings.ToLower(command.Name) == strings.ToLower(command_name) {
            if command.Admin {
                if CheckAdmin(m.Author.ID) {
                    command.Command(s, m, args)
                } else {
                    s.ChannelMessageSend(support.Config.FactorioChannelID, "You are not an admin!")
                }
            }

            if !command.Admin {
                command.Command(s, m, args)
            }
            return
        }
    }
    s.ChannelMessageSend(support.Config.FactorioChannelID, "Command not found!")
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
