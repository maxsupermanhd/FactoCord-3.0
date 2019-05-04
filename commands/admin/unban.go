package admin

import (
	"io"

	"../../support"
	"github.com/bwmarrin/discordgo"
)

// SaveServer executes the save command on the server.
func UnbanPlayer(s *discordgo.Session, m *discordgo.MessageCreate, arg1 string) {
	io.WriteString(*P, "/unban " + arg1 + "\n")
	s.ChannelMessageSend(support.Config.FactorioChannelID, "Player "+ arg1 + " unbanned!")
	return
}
