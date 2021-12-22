package admin

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

var BanPlayerDoc = support.Command{
	Name:  "ban",
	Desc:  "Ban a user from the server",
	Usage: "/ban <player> <reason>",
	Doc:   `command bans the player on the server with a specified reason`,
	Admin: true,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "player",
			Description: "In-game nick name of the player to ban",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "reason",
			Description: "Reason for banning",
			Required:    true,
		},
	},
	Command: BanPlayer,
}

// BanPlayer bans a player on the server.
func BanPlayer(s *discordgo.Session, i *discordgo.InteractionCreate) {
	player := strings.TrimSpace(i.ApplicationCommandData().Options[0].StringValue())
	reason := strings.TrimSpace(i.ApplicationCommandData().Options[1].StringValue())

	command := "/ban " + player + " " + reason
	success := support.Factorio.Send(command)
	if success {
		support.Respond(s, i, "Player "+player+" banned with reason \""+reason+"\"!")
	} else {
		support.Respond(s, i, "Sorry, there was an error sending /ban command")
	}
}
