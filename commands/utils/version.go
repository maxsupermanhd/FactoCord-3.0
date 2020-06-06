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
	factorioVersion := strings.Fields(string(out))[1]
	if err != nil {
		s.ChannelMessageSend(support.Config.FactorioChannelID, fmt.Sprintf("Sorry, there was an error. Error details: %s", err))
	}
	res := "Server version: **" + factorioVersion + "**"

	cmd = exec.Command("git", "describe", "--tags")
	out, err = cmd.CombinedOutput()
	if err != nil {
		s.ChannelMessageSend(support.Config.FactorioChannelID, fmt.Sprintf("Sorry, there was an error. Error details: %s", err))
	}
	factocord := "FactoCord version: **" + string(out) + "**"
	res += "\n" + factocord

	s.ChannelMessageSend(support.Config.FactorioChannelID, res)
	return
}
