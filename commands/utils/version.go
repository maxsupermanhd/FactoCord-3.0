package utils

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/maxsupermanhd/FactoCord-3.0/v3/support"
)

var VersionDoc = support.CommandDoc{
	Name: "version",
	Doc: `command outputs factorio server version and FactoCord version.
If it says that FactoCord version is unknown look into the error.log`,
}

func VersionString(s *discordgo.Session, _ string) {
	factorioVersion, err := support.FactorioVersion()
	if err != nil {
		support.Send(s, "Sorry, there was an error checking factorio version")
		support.Panik(err, "... when running `factorio --version`")
		return
	}
	res := "Server version: **" + factorioVersion + "**"

	res += fmt.Sprintf("\nFactoCord version: **%s**", support.FactoCordVersion)

	support.Send(s, res)
}
