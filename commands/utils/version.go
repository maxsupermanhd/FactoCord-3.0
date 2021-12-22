package utils

import (
	"fmt"
	"github.com/bwmarrin/discordgo"

	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

var VersionDoc = support.Command{
	Name: "version",
	Desc: "Get server version",
	Doc: `command outputs factorio server version and FactoCord version.
If it says that FactoCord version is unknown look into the error.log`,
	Command: VersionString,
}

func VersionString(s *discordgo.Session, i *discordgo.InteractionCreate) {
	factorioVersion, err := support.FactorioVersion()
	if err != nil {
		support.Respond(s, i, "Sorry, there was an error checking factorio version")
		support.Panik(err, "... when running `factorio --version`")
		return
	}
	res := "Server version: **" + factorioVersion + "**"

	res += fmt.Sprintf("\nFactoCord version: **%s**", support.FactoCordVersion)

	support.Respond(s, i, res)
}
