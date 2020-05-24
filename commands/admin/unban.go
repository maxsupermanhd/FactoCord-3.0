package admin

import (
    "io"
    "strings"

    "../../support"
    "github.com/bwmarrin/discordgo"
)

var UnbanPlayerUsage = "Usage $unban <player>"

// SaveServer executes the save command on the server.
func UnbanPlayer(s *discordgo.Session, m *discordgo.MessageCreate, args string) {
    if strings.ContainsAny(args, " \n\t") {
        s.ChannelMessageSend(support.Config.FactorioChannelID, support.FormatUsage(UnbanPlayerUsage))
        return
    }
    io.WriteString(*P, "/unban " + args + "\n")
    s.ChannelMessageSend(support.Config.FactorioChannelID, "Player "+ args + " unbanned!")
    return
}
