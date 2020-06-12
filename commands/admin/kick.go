package admin

import (
	"io"
	"strings"

	"../../support"
	"github.com/bwmarrin/discordgo"
)

// KickPlayerUsage comment...
var KickPlayerUsage = "Usage: $kick <player> <reason>"

// KickPlayer kicks a player from the server.
func KickPlayer(s *discordgo.Session, m *discordgo.MessageCreate, args string) {
	if len(args) == 0 {
		support.SendFormat(s, KickPlayerUsage)
		return
	}
	args2 := strings.SplitN(args+" ", " ", 2)
	player := strings.TrimSpace(args2[0])
	reason := strings.TrimSpace(args2[1])

	if len(player) == 0 || len(reason) == 0 {
		support.SendFormat(s, KickPlayerUsage)
		return
	}
	command := "/kick " + player + " " + reason + "\n"
	_, err := io.WriteString(*P, command)
	if err != nil {
		support.Send(s, "Sorry, there was an error sending /kick command")
		support.Panik(err, "... when sending \""+command+"\"")
		return
	}
	support.Send(s, "Player "+player+" kicked with reason "+reason+"!")
	return
}
