package admin

import (
	"io"

	"../../support"
	"github.com/bwmarrin/discordgo"
)

// SaveServer executes the save command on the server.
func SaveServer(s *discordgo.Session, m *discordgo.MessageCreate, args string) {
	if len(args) != 0 {
		support.Send(s, "Save accepts no arguments")
		return
	}
	_, err := io.WriteString(*P, "/save\n")
	if err != nil {
		support.Send(s, "Sorry, there was an error sending /save command")
		support.Panik(err, "... when sending \"/save\"")
		return
	}
	support.Send(s, "Server saved successfully!")
	return
}
