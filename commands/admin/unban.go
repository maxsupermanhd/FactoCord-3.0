package admin

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

var UnbanPlayerDoc = support.CommandDoc{
	Name:  "unban",
	Usage: "$unban <player>",
	Doc:   `command removes the player from the banlist on the server`,
}

// UnbanPlayer unbans a player on the server.
func UnbanPlayer(s *discordgo.Session, args string) {
	if strings.ContainsAny(args, " \n\t") {
		support.SendFormat(s, "Usage: "+UnbanPlayerDoc.Usage)
		return
	}
	command := "/unban " + args
	success := support.Factorio.Send(command)
	if success {
		support.Send(s, "Player "+args+" unbanned!")
	} else {
		support.Send(s, "Sorry, there was an error sending /unban command")
	}
}
