package admin

import (
	"io"

	"../../support"
	"github.com/bwmarrin/discordgo"
)

// SaveServer executes the save command on the server.
func BanPlayer(s *discordgo.Session, m *discordgo.MessageCreate, arg1 string, arg2 string) {
	io.WriteString(*P, "/ban " + arg1 + " " + arg2 +"\n")
	s.ChannelMessageSend(support.Config.FactorioChannelID, "Player "+ arg1 + " banned with reason " + arg2 + "!")
	return
}
