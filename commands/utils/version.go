package utils

import (
	"os/exec"
	"strings"

	"../../support"
	"github.com/bwmarrin/discordgo"
)

var VersionStringUsage = "Usage: $version"

func VersionString(s *discordgo.Session, _ string) {
	cmd := exec.Command(support.Config.Executable, "--version")
	out, err := cmd.CombinedOutput()
	factorioVersion := strings.Fields(string(out))[1]
	if err != nil {
		support.Send(s, "Sorry, there was an error checking factorio version")
		support.Panik(err, "... when running `factorio --version`")
		return
	}
	res := "Server version: **" + factorioVersion + "**"

	cmd = exec.Command("git", "describe", "--tags")
	out, err = cmd.CombinedOutput()
	if err != nil {
		support.Send(s, "Sorry, there was an error checking git version")
		support.Panik(err, "... when running `git describe --tags`")
		return
	}
	factocord := "FactoCord version: **" + string(out) + "**"
	res += "\n" + factocord

	support.Send(s, res)
}
