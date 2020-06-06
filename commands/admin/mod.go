package admin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"../../support"
	"github.com/bwmarrin/discordgo"
)

// ModJSON is struct containing a slice of Mod.
type ModJSON struct {
	Mods []Mod `json:"mods"`
}

// Mod is a struct containing info about a mod.
type Mod struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Version string `json:"version,omitempty"`
}

// ModCommandUsage ...
var ModCommandUsage = "Usage: $mod purge | (add|remove|enable|disable) <modnames>"

// ModCommand returns the list of mods running on the server.
func ModCommand(s *discordgo.Session, m *discordgo.MessageCreate, args string) {
	argsList := strings.Split(args, " ")
	if len(argsList) == 0 {
		s.ChannelMessageSend(support.Config.FactorioChannelID, support.FormatUsage(ModCommandUsage))
		return
	}

	action := argsList[0]
	if action == "add" || action == "remove" || action == "enable" || action == "disable" {
		if len(argsList) < 2 {
			s.ChannelMessageSend(support.Config.FactorioChannelID, support.FormatUsage("Usage: $mod "+action+" <modname> [<modname>]+"))
			return
		}
	} else if action != "purge" {
		s.ChannelMessageSend(support.Config.FactorioChannelID, support.FormatUsage(ModCommandUsage))
		return
	}

	mods := &ModJSON{}
	modsListFile, err := ioutil.ReadFile(support.Config.ModListLocation)
	// Don't exit on this error, just sent message to the channel!
	if err != nil {
		s.ChannelMessageSend(support.Config.FactorioChannelID,
			fmt.Sprintf("Sorry, there was an error reading your mods list, did you specify it in the .env file? Error details: %s", err))
		return
	}

	err = json.Unmarshal(modsListFile, &mods)
	if err != nil {
		s.ChannelMessageSend(support.Config.FactorioChannelID,
			fmt.Sprintf("Sorry, there was an error reading your mods list. Error details: %s", err))
		return
	}

	modnames := support.DeleteEmptyStrings(argsList[1:])
	var res string
	if action == "purge" {
		res = modsPurge(mods)
	} else if action == "add" {
		res = modsAdd(mods, modnames)
	} else if action == "remove" {
		res = modsRemove(mods, modnames)
	} else if action == "enable" {
		res = modsEnable(mods, modnames, true)
	} else if action == "disable" {
		res = modsEnable(mods, modnames, false)
	}

	modsListFile, err = json.MarshalIndent(mods, "", "    ")
	if err != nil {
		s.ChannelMessageSend(support.Config.FactorioChannelID,
			fmt.Sprintf("Sorry, there was an error. Error details: %s", err))
		return
	}
	ioutil.WriteFile(support.Config.ModListLocation, modsListFile, 0666)

	support.ChunkedMessageSend(s, support.Config.FactorioChannelID, res)
	return
}

func modsPurge(mods *ModJSON) string {
	total := 0
	purged := 0
	end := 0
	res := ""
	for i, mod := range mods.Mods {
		if mod.Enabled {
			mods.Mods[end] = mods.Mods[i]
			end++
		} else {
			purged++
			res += "\n    " + mod.Name
		}
		total++
	}
	mods.Mods = mods.Mods[:end]
	res = fmt.Sprintf("**Removed %d disabled mods** (left: %d, total: %d):", purged, len(mods.Mods), total) + res
	return res
}

func modsAdd(mods *ModJSON, modnames []string) string {
	modsList := make([]Mod, len(mods.Mods)+len(modnames))

	end := len(mods.Mods)
	copy(modsList, mods.Mods)
	mods.Mods = modsList

	res := ""
	alreadyAdded := ""

	for _, modname := range modnames {
		added := false
		for i := 0; i < end; i++ {
			mod := mods.Mods[i]
			if strings.ToLower(mod.Name) == strings.ToLower(modname) {
				alreadyAdded += "\n    " + modname
				added = true
				break
			}
			if strings.ToLower(mod.Name) > strings.ToLower(modname) {
				copy(mods.Mods[i+1:], mods.Mods[i:])
				mods.Mods[i] = Mod{
					Name:    modname,
					Enabled: true,
				}
				end++
				added = true
				res += "\n    " + modname
				break
			}
		}
		if !added {
			res += "\n    " + modname
			mods.Mods[end] = Mod{
				Name:    modname,
				Enabled: true,
			}
			end++
		}
	}
	if len(modnames) == 1 {
		if len(alreadyAdded) > 0 {
			res = "Mod \"" + modnames[0] + "\" is already added"
		} else {
			res = "Added mod \"" + modnames[0] + "\""
		}
	} else {
		res = "**Added mods:**" + res
		if len(alreadyAdded) > 0 {
			res += "\n**Already added:**" + alreadyAdded
		}
	}
	mods.Mods = mods.Mods[:end]
	return res
}

func modsRemove(mods *ModJSON, modnames []string) string {
	removed := 0
	res := ""
	notFoundCount := 0
	notFound := ""
	for _, modname := range modnames {
		found := false

		for i, mod := range mods.Mods {
			if modname == mod.Name {
				found = true
				res += "\n    " + modname
				copy(mods.Mods[i:], mods.Mods[i+1:])
				removed++
				break
			}
		}
		if !found {
			notFoundCount++
			notFound += "\n    " + modname
		}
	}
	mods.Mods = mods.Mods[:len(mods.Mods)-removed]
	if len(modnames) == 1 {
		if notFoundCount > 0 {
			res = "Mod \"" + modnames[0] + "\" not found"
		} else {
			res = "Removed mod \"" + modnames[0] + "\""
		}
	} else {
		res = fmt.Sprintf("**Removed %d mods (left: %d):**", removed, len(mods.Mods)) + res
		if notFoundCount > 0 {
			notFound = fmt.Sprintf("\n**Not Found %d mods:**", notFoundCount) + notFound
			res += notFound
		}
	}

	return res
}

func modsEnable(mods *ModJSON, modnames []string, setEnabled bool) string {
	res := ""
	notFound := ""
	notFoundCount := 0

	count := 0
	for _, modname := range modnames {
		found := false
		for i, mod := range mods.Mods {
			if mod.Name == modname {
				mods.Mods[i].Enabled = setEnabled
				found = true
				count++
				res += "\n    " + modname
			}
		}
		if !found {
			notFoundCount++
			notFound += "\n    " + modname
		}
	}

	action := "Disabled"
	if setEnabled {
		action = "Enabled"
	}
	if len(modnames) == 1 {
		if len(notFound) > 0 {
			res = "Mod \"" + modnames[0] + "\" not found"
		} else {
			res = action + " mod \"" + modnames[0] + "\""
		}
	} else {
		res = fmt.Sprintf("**"+action+" %d mods:**", count) + res
		if len(notFound) > 0 {
			notFound = fmt.Sprintf("\n**Not Found %d mods:**", notFoundCount) + notFound
			res += notFound
		}
	}
	return res
}
