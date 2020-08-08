package admin

import (
	"github.com/bwmarrin/discordgo"

	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

var SaveServerDoc = support.CommandDoc{
	Name: "save",
	Doc:  `command sends a command to save the game to the server`,
}

// SaveServer executes the save command on the server.
func SaveServer(s *discordgo.Session, args string) {
	if len(args) != 0 {
		support.Send(s, "Save accepts no arguments")
		return
	}
	success := support.Factorio.Send("/save")
	if success {
		support.Factorio.SaveRequested = true
		//support.Send(s, "Server saved successfully!")
	} else {
		support.Send(s, "Sorry, there was an error sending /save command")
	}
}
