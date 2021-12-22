package admin

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

// KickPlayerUsage comment...

var KickPlayerDoc = support.Command{
	Name:  "kick",
	Desc:  "Kick a user from the server",
	Usage: "/kick <player> <reason>",
	Doc:   `command kicks the player out from the server with a specified reason`,
	Admin: true,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "player",
			Description: "In-game nick name of the player to kick",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "reason",
			Description: "Reason for kicking",
			Required:    true,
		},
	},
	Command: KickPlayer,
}

// KickPlayer kicks a player from the server.
func KickPlayer(s *discordgo.Session, i *discordgo.InteractionCreate) {
	player := strings.TrimSpace(i.ApplicationCommandData().Options[0].StringValue())
	reason := strings.TrimSpace(i.ApplicationCommandData().Options[1].StringValue())

	command := "/kick " + player + " " + reason
	success := support.Factorio.Send(command)
	if success {
		support.Respond(s, i, "Player "+player+" kicked with reason "+reason+"!")
	} else {
		support.Respond(s, i, "Sorry, there was an error sending /kick command")
	}
}
