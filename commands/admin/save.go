package admin

import (
    "io"

    "../../support"
    "github.com/bwmarrin/discordgo"
)

// SaveServer executes the save command on the server.
func SaveServer(s *discordgo.Session, m *discordgo.MessageCreate, args string) {
    if len(args) != 0 {
        s.ChannelMessageSend(support.Config.FactorioChannelID, "Save accepts no arguments")
        return
    }
    io.WriteString(*P, "/save\n")
    s.ChannelMessageSend(support.Config.FactorioChannelID, "Server saved successfully!")
    return
}
