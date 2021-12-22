package admin

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

var SaveServerDoc = support.Command{
	Name:    "save",
	Desc:    "Save the game",
	Doc:     `command sends a command to save the game to the server`,
	Admin:   true,
	Command: SaveServer,
}

// SaveServer executes the save command on the server.
func SaveServer(s *discordgo.Session, i *discordgo.InteractionCreate) {
	success := support.Factorio.Send("/save")
	if success {
		support.Factorio.SaveRequested = true
		support.Respond(s, i, "Server is saving the game")
	} else {
		support.Respond(s, i, "Sorry, there was an error sending /save command")
	}
}
