package admin

import (
	"io"
	"time"

	"../../support"
	"github.com/bwmarrin/discordgo"
)

// R is a reference to the running boolean in main.
var R *bool

// RestartCount is the number of times the server has restarted.
var RestartCount int

// Restart saves and restarts the server
func Restart(s *discordgo.Session, m *discordgo.MessageCreate, args string) {
	if len(args) != 0 {
		support.Send(s, "Restart accepts no arguments")
		return
	}
	if *R == false {
		support.Send(s, "Server is not running!")
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
	support.Send(s, "Saved server, now restarting!")
	time.Sleep(3 * time.Second)
	*R = false
	RestartCount = RestartCount + 1
	return
}
