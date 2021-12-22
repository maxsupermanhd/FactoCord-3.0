package utils

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"path"

	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

// ModJson is struct containing a slice of Mod.
type ModJson struct {
	Mods []Mod
}

// Mod is a struct containing info about a mod.
type Mod struct {
	Name    string
	Enabled bool
}

var ModListDoc = support.Command{
	Name:    "mods",
	Desc:    "List the mods on the server",
	Doc:     `command outputs information about current mods`,
	Usage:   "/mods on | off | all | files",
	Command: ModsList,
	Subcommands: []support.Command{
		{
			Name: "on",
			Desc: `Shows currently enabled mods`,
			Doc:  `command sends currently enabled mods`,
		},
		{
			Name: "off",
			Desc: `Shows currently disabled mods`,
			Doc:  `command sends currently disabled mods`,
		},
		{
			Name: "all",
			Desc: `Shows all mods in mod-list.json`,
			Doc:  `command sends all mods in mod-list.json`,
		},
		{
			Name: "files",
			Desc: `Shows filenames of all downloaded mods`,
			Doc:  `command sends filenames of all downloaded mods`,
		},
	},
}

func modList(ModList *ModJson, returnEnabled bool, returnDisabled bool) string {
	var enabled, disabled int
	var S = "mod"
	if len(ModList.Mods) > 1 {
		S = "mods"
	}
	for _, mod := range ModList.Mods {
		if mod.Enabled {
			enabled += 1
		} else {
			disabled += 1
		}
	}

	res := fmt.Sprintf("%d total %s (%d enabled, %d disabled)", len(ModList.Mods), S, enabled, disabled)

	if returnEnabled {
		res += "\n**Enabled:**"
		any := false
		for _, mod := range ModList.Mods {
			if mod.Enabled {
				any = true
				res += "\n    " + mod.Name
			}
		}
		if !any {
			res += " **None**"
		}
	}
	if returnDisabled {
		if returnEnabled {
			res += "\n"
		}
		res += "\n**Disabled:**"
		any := false
		for _, mod := range ModList.Mods {
			if !mod.Enabled {
				any = true
				res += "\n    " + mod.Name
			}
		}
		if !any {
			res += " **None**"
		}
	}

	return res
}

func modsFiles() string {
	res := ""
	baseDir := path.Dir(support.Config.ModListLocation)
	files, err := ioutil.ReadDir(baseDir)
	if err != nil {
		support.Critical(err, "wtf")
	}
	for _, file := range files {
		re := support.ModFileRegexp.FindString(file.Name())
		if re != "" {
			res += "\n    " + file.Name()
		}
	}
	if res == "" {
		return "**No mods**"
	} else {
		return "**Installed mods:**" + res
	}
}

// ModsList returns the list of mods running on the server.
func ModsList(s *discordgo.Session, i *discordgo.InteractionCreate) {
	returnEnabled := true
	returnDisabled := false
	command := i.ApplicationCommandData().Options[0].Name
	if command == "on" || command == "" {
		returnEnabled = true
	} else if command == "off" {
		returnEnabled = false
		returnDisabled = true
	} else if command == "all" {
		returnDisabled = true
	} else if command == "files" {
		support.Respond(s, i, modsFiles())
		return
	}
	ModList := &ModJson{}
	Json, err := ioutil.ReadFile(support.Config.ModListLocation)
	if err != nil {
		support.Respond(s, i, "Sorry, there was an error reading your mods list")
		support.Panik(err, "there was an error reading mods list, did you specify it in the config.json file?")
		return
	}

	err = json.Unmarshal(Json, &ModList)
	if err != nil {
		support.Respond(s, i, "Sorry, there was an error reading your mods list")
		support.Panik(err, "there was an error reading mods list")
		return
	}

	support.RespondChunked(s, i, modList(ModList, returnEnabled, returnDisabled))
	return
}
