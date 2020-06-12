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
		support.Send(s, "Stop accepts no arguments")
		return
	}
	_, err := io.WriteString(*P, "/save\n")
	if err != nil {
		support.Send(s, "Sorry, there was an error sending /save command")
		support.Panik(err, "... when sending \"/save\"")
		return
	}
	_, err = io.WriteString(*P, "/quit\n")
	if err != nil {
		support.Send(s, "Sorry, there was an error sending /quit command")
		support.Panik(err, "... when sending \"/quit\"")
		return
	}
	support.Send(s, "Server saved and shutting down; Cya!")
	err = s.Close()
	support.Critical(err, "... when closing discord connection")
	support.Exit(0)
}
