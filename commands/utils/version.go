package utils

import (
	"fmt"
	"os/exec"
	"strings"
	
	"../../support"
	"github.com/bwmarrin/discordgo"
)

var VersionStringUsage = "Usage: $version"

func VersionString(s *discordgo.Session, m *discordgo.MessageCreate, ar string) {
	cmd := exec.Command(support.Config.Executable, "--version")
    out, err := cmd.CombinedOutput()
	if err != nil {
		s.ChannelMessageSend(support.Config.FactorioChannelID, fmt.Sprintf("Sorry, there was an error. Error details: %s", err))
	}
	s.ChannelMessageSend(support.Config.FactorioChannelID, "Server version: " + strings.Fields(string(out))[1]);
	return
}


