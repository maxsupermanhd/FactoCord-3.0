package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"../../support"
	"github.com/bwmarrin/discordgo"
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

var ModListUsage = "Usage: $mods [on | off | all]"

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
		for _, mod := range ModList.Mods {
			if mod.Enabled {
				res += "\n    " + mod.Name
			}
		}
	}
	if returnDisabled {
		if returnEnabled {
			res += "\n"
		}
		res += "\n**Disabled:**"
		for _, mod := range ModList.Mods {
			if !mod.Enabled {
				res += "\n    " + mod.Name
			}
		}
	}

	return res
}

// ModsList returns the list of mods running on the server.
func ModsList(s *discordgo.Session, args string) {
	returnEnabled := true
	returnDisabled := false
	if args == "on" || args == "" {
		returnEnabled = true
	} else if args == "off" {
		returnEnabled = false
		returnDisabled = true
	} else if args == "all" {
		returnDisabled = true
	} else {
		support.SendFormat(s, ModListUsage)
		return
	}
	ModList := &ModJson{}
	Json, err := ioutil.ReadFile(support.Config.ModListLocation)
	if err != nil {
		support.Send(s, "Sorry, there was an error reading your mods list")
		support.Panik(err, "there was an error reading mods list, did you specify it in the .env file?")
		return
	}

	err = json.Unmarshal(Json, &ModList)
	if err != nil {
		support.Send(s, "Sorry, there was an error reading your mods list")
		support.Panik(err, "there was an error reading mods list")
		return
	}
	support.ChunkedMessageSend(s, modList(ModList, returnEnabled, returnDisabled))
	return
}
