package admin

import (
	"../../support"
	"github.com/bwmarrin/discordgo"
)

// SaveServer executes the save command on the server.
func SaveServer(s *discordgo.Session, m *discordgo.MessageCreate, args string) {
	if len(args) != 0 {
		support.Send(s, "Save accepts no arguments")
		return
	}
	success := support.SendToFactorio("/save")
	if success {
		// TODO read log to be sure it's successful
		support.Send(s, "Server saved successfully!")
	} else {
		support.Send(s, "Sorry, there was an error sending /save command")
	}
}
