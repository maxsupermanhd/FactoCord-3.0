package admin

import (
	"../../support"
	"github.com/bwmarrin/discordgo"
)

// StopServer saves and stops the server.
func StopServer(s *discordgo.Session, args string) {
	if len(args) != 0 {
		support.Send(s, "Stop accepts no arguments")
		return
	}
	success := support.SendToFactorio("/save")
	if !success {
		support.Send(s, "Sorry, there was an error sending /save command")
		return
	}
	success = support.SendToFactorio("/quit")
	if !success {
		support.Send(s, "Sorry, there was an error sending /quit command")
		return
	}
	//support.Send(s, "Server saved and shutting down; Cya!")
	err := s.Close()
	support.Critical(err, "... when closing discord connection")
	support.Exit(0)
}
