package admin

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

// UnbanPlayerUsage comment...
var UnbanPlayerUsage = "Usage $unban <player>"

// UnbanPlayer unbans a player on the server.
func UnbanPlayer(s *discordgo.Session, args string) {
	if strings.ContainsAny(args, " \n\t") {
		support.SendFormat(s, UnbanPlayerUsage)
		return
	}
	command := "/unban " + args
	success := support.SendToFactorio(command)
	if success {
		support.Send(s, "Player "+args+" unbanned!")
	} else {
		support.Send(s, "Sorry, there was an error sending /unban command")
	}
}
