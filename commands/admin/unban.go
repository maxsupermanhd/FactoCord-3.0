package admin

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

var UnbanPlayerDoc = support.Command{
	Name:  "unban",
	Desc:  "Unban a user from the server",
	Usage: "/unban <player>",
	Doc:   `command removes the player from the banlist on the server`,
	Admin: true,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "player",
			Description: "In-game nick name of the player to unban",
			Required:    true,
		},
	},
	Command: UnbanPlayer,
}

// UnbanPlayer unbans a player on the server.
func UnbanPlayer(s *discordgo.Session, i *discordgo.InteractionCreate) {
	player := strings.TrimSpace(i.ApplicationCommandData().Options[0].StringValue())

	command := "/unban " + player
	success := support.Factorio.Send(command)
	if success {
		support.Send(s, "Player "+player+" unbanned!")
	} else {
		support.Send(s, "Sorry, there was an error sending /unban command")
	}
}
