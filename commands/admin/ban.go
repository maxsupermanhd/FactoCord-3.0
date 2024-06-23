package admin

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/maxsupermanhd/FactoCord-3.0/v3/support"
)

var BanPlayerDoc = support.CommandDoc{
	Name:  "ban",
	Usage: "$ban <player> <reason>",
	Doc:   `command bans the player on the server with a specified reason`,
}

// BanPlayer bans a player on the server.
func BanPlayer(s *discordgo.Session, args string) {
	if len(args) == 0 {
		support.SendFormat(s, "Usage: "+BanPlayerDoc.Usage)
		return
	}
	args2 := strings.SplitN(args+" ", " ", 2)
	player := strings.TrimSpace(args2[0])
	reason := strings.TrimSpace(args2[1])

	if len(player) == 0 || len(reason) == 0 {
		support.SendFormat(s, "Usage: "+BanPlayerDoc.Usage)
		return
	}

	command := "/ban " + player + " " + reason
	success := support.Factorio.Send(command)
	if success {
		support.Send(s, "Player "+player+" banned with reason \""+reason+"\"!")
	} else {
		support.Send(s, "Sorry, there was an error sending /ban command")
	}
}
