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

func modList(ModList *ModJson, return_enabled bool, return_disabled bool) string {
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

	if return_enabled {
		res += "\n**Enabled:**"
		for _, mod := range ModList.Mods {
			if mod.Enabled {
				res += "\n    " + mod.Name
			}
		}
	}
	if return_disabled {
		if return_enabled {
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
func ModsList(s *discordgo.Session, m *discordgo.MessageCreate, args string) {
	return_enabled := true
	return_disabled := false
	if args == "on" || args == "" {
		return_enabled = true
	} else if args == "off" {
		return_enabled = false
		return_disabled = true
	} else if args == "all" {
		return_disabled = true
	} else {
		s.ChannelMessageSend(support.Config.FactorioChannelID, support.FormatUsage(ModListUsage))
		return
	}
	ModList := &ModJson{}
	Json, err := ioutil.ReadFile(support.Config.ModListLocation)
	// Don't exit on this error, just sent message to the channel!
	if err != nil {
		s.ChannelMessageSend(support.Config.FactorioChannelID,
			fmt.Sprintf("Sorry, there was an error reading your mods list, did you specify it in the .env file? Error details: %s", err))
		return
	}

	err = json.Unmarshal(Json, &ModList)
	if err != nil {
		s.ChannelMessageSend(support.Config.FactorioChannelID,
			fmt.Sprintf("Sorry, there was an error reading your mods list. Error details: %s", err))
		return
	}
	support.ChunkedMessageSend(s, support.Config.FactorioChannelID, modList(ModList, return_enabled, return_disabled))
	return
}
