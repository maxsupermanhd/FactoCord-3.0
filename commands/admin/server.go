package admin

import (
	"../../support"
	"github.com/bwmarrin/discordgo"
)

var ServerCommandUsage = "Usage: $server [stop|start|restart]"

func ServerCommand(s *discordgo.Session, args string) {
	switch args {
	case "":
		if support.Factorio.IsRunning() {
			support.Send(s, "Factorio server is **running**")
		} else {
			support.Send(s, "Factorio server is **stopped**")
		}
	case "stop":
		support.Factorio.Stop(s)
	case "start":
		support.Factorio.Start(s)
	case "restart":
		support.Factorio.Stop(s)
		support.Factorio.Start(s)
	default:
		support.SendFormat(s, ServerCommandUsage)
	}
}
