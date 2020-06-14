package admin

import (
	"strings"

	"../../support"
	"github.com/bwmarrin/discordgo"
)

// BanPlayerUsage comment...
var BanPlayerUsage = "Usage: $ban <player> <reason>"

// BanPlayer bans a player on the server.
func BanPlayer(s *discordgo.Session, args string) {
	if len(args) == 0 {
		support.SendFormat(s, BanPlayerUsage)
		return
	}
	args2 := strings.SplitN(args+" ", " ", 2)
	player := strings.TrimSpace(args2[0])
	reason := strings.TrimSpace(args2[1])

	if len(player) == 0 || len(reason) == 0 {
		support.SendFormat(s, BanPlayerUsage)
		return
	}

	command := "/ban " + player + " " + reason
	success := support.SendToFactorio(command)
	if success {
		support.Send(s, "Player "+player+" banned with reason \""+reason+"\"!")
	} else {
		support.Send(s, "Sorry, there was an error sending /ban command")
	}
}
