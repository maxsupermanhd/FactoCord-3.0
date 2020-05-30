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
		s.ChannelMessageSend(support.Config.FactorioChannelID, support.FormatUsage(KickPlayerUsage))
		return
	}
	args2 := strings.SplitN(args + " ", " ", 2)
	player := strings.TrimSpace(args2[0])
	reason := strings.TrimSpace(args2[1])

	if len(player) == 0 || len(reason) == 0 {
		s.ChannelMessageSend(support.Config.FactorioChannelID, support.FormatUsage(KickPlayerUsage))
		return
	}
	io.WriteString(*P, "/kick " + player + " " + reason + "\n")
	s.ChannelMessageSend(support.Config.FactorioChannelID, "Player "+ player + " kicked with reason " + reason + "!")
	return
}
