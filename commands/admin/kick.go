package admin

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

// KickPlayerUsage comment...

var KickPlayerDoc = support.CommandDoc{
	Name:  "kick",
	Usage: "$kick <player> <reason>",
	Doc:   `command kicks the player out from the server with a specified reason`,
}

// KickPlayer kicks a player from the server.
func KickPlayer(s *discordgo.Session, args string) {
	if len(args) == 0 {
		support.SendFormat(s, "Usage: "+KickPlayerDoc.Usage)
		return
	}
	args2 := strings.SplitN(args+" ", " ", 2)
	player := strings.TrimSpace(args2[0])
	reason := strings.TrimSpace(args2[1])

	if len(player) == 0 || len(reason) == 0 {
		support.SendFormat(s, "Usage: "+KickPlayerDoc.Usage)
		return
	}
	command := "/kick " + player + " " + reason
	success := support.Factorio.Send(command)
	if success {
		support.Send(s, "Player "+player+" kicked with reason "+reason+"!")
	} else {
		support.Send(s, "Sorry, there was an error sending /kick command")
	}
}
