package admin

import (
	"io"

	"../../support"
	"github.com/bwmarrin/discordgo"
)

// SaveServer executes the save command on the server.
func KickPlayer(s *discordgo.Session, m *discordgo.MessageCreate, arg1 string, arg2 string) {
	io.WriteString(*P, "/kick " + arg1 + " " + arg2 + "\n")
	s.ChannelMessageSend(support.Config.FactorioChannelID, "Player "+ arg1 + " kicked with reason " + arg2 + "!")
	return
}
