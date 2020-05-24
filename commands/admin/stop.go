package admin

import (
    "io"

    "../../support"
    "github.com/bwmarrin/discordgo"
)

// P references the var Pipe in main
var P *io.WriteCloser

// StopServer saves and stops the server.
func StopServer(s *discordgo.Session, m *discordgo.MessageCreate, args string) {
    if len(args) != 0 {
        s.ChannelMessageSend(support.Config.FactorioChannelID, "Stop accepts no arguments")
        return
    }
    io.WriteString(*P, "/save\n")
    io.WriteString(*P, "/quit\n")
    s.ChannelMessageSend(support.Config.FactorioChannelID, "Server saved and shutting down; Cya!")
    s.Close()
    support.Exit(0)
}
