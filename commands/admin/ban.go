package admin

import (
	"io"
	"strings"

	"../../support"
	"github.com/bwmarrin/discordgo"
)

var BanPlayerUsage = "Usage: $ban <player> <reason>"

// SaveServer executes the save command on the server.
func BanPlayer(s *discordgo.Session, m *discordgo.MessageCreate, args string) {
	if len(args) == 0 {
		s.ChannelMessageSend(support.Config.FactorioChannelID, support.FormatUsage(BanPlayerUsage))
		return
	}
	args2 := strings.SplitN(args + " ", " ", 2)
	player := strings.TrimSpace(args2[0])
	reason := strings.TrimSpace(args2[1])

	if len(player) == 0 || len(reason) == 0 {
		s.ChannelMessageSend(support.Config.FactorioChannelID, support.FormatUsage(BanPlayerUsage))
		return
	}

	io.WriteString(*P, "/ban " + player + " " + reason +"\n")
	s.ChannelMessageSend(support.Config.FactorioChannelID, "Player "+ player + " banned with reason \"" + reason + "\"!")
	return
}
