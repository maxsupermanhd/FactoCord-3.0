package utils

import (
	"io/ioutil"
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

	factocord := "FactoCord version unknown"
	if support.DirExists("./.git") {
		gitNotInstalledErr := exec.Command("sh", "-c", "command -v git").Run()
		cmd = exec.Command("git", "describe", "--tags")
		out, err = cmd.CombinedOutput()
		if err != nil {
			if gitNotInstalledErr != nil {
				factocord += ": git is probably not installed"
				support.Panik(gitNotInstalledErr, "Fail running `sh -c 'command -v git'` to check if git is installed")
			} else {
				support.Send(s, "Sorry, there was an error checking git version")
				support.Panik(err, "... when running `git describe --tags`")
				return
			}
		}
		factocord = "FactoCord version: **" + string(out) + "**"
	} else if support.FileExists("./.version") {
		version, err := ioutil.ReadFile("./.version")
		if err == nil {
			factocord = "FactoCord version: **" + strings.TrimSpace(string(version)) + "**"
		} else {
			support.Panik(err, "... when reading .version")
		}
	}
	res += "\n" + factocord

	support.Send(s, res)
}
