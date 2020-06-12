package admin

import (
	"io"
	"strings"

	"../../support"
	"github.com/bwmarrin/discordgo"
)

// UnbanPlayerUsage comment...
var UnbanPlayerUsage = "Usage $unban <player>"

// UnbanPlayer unbans a player on the server.
func UnbanPlayer(s *discordgo.Session, m *discordgo.MessageCreate, args string) {
	if strings.ContainsAny(args, " \n\t") {
		support.SendFormat(s, UnbanPlayerUsage)
		return
	}
	command := "/unban " + args + "\n"
	_, err := io.WriteString(*P, command)
	if err != nil {
		support.Send(s, "Sorry, there was an error sending /unban command")
		support.Panik(err, "... when sending \""+command+"\"")
		return
	}
	support.Send(s, "Player "+args+" unbanned!")
	return
}
