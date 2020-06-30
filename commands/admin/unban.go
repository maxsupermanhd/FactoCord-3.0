package admin

import (
	"strings"

	"../../support"
	"github.com/bwmarrin/discordgo"
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
	success := support.Factorio.Send(command)
	if success {
		support.Send(s, "Player "+args+" unbanned!")
	} else {
		support.Send(s, "Sorry, there was an error sending /unban command")
	}
}
