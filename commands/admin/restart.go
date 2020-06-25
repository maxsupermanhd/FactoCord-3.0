package admin

import (
	"time"

	"../../support"
	"github.com/bwmarrin/discordgo"
)

// R is a reference to the running boolean in main.
var R *bool

// RestartCount is the number of times the server has restarted.
var RestartCount int

// Restart saves and restarts the server
func Restart(s *discordgo.Session, args string) {
	if len(args) != 0 {
		support.Send(s, "Restart accepts no arguments")
		return
	}
	if *R == false {
		support.Send(s, "Server is not running!")
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
	//support.Send(s, "Saved server, now restarting!")
	// TODO wait for factorio to exit
	time.Sleep(3 * time.Second)
	*R = false
	RestartCount = RestartCount + 1
	return
}
