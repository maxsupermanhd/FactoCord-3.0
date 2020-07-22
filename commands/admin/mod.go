package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/maxsupermanhd/FactoCord-3.0/support"
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

type modDescriptionT struct {
	name    string
	path    string
	version support.SemanticVersionT
}

func modDescription(s string) (*modDescriptionT, *error) {
	name, version := support.SplitDivide(s, "==")
	version2, err := support.SemanticVersion(version)
	if err != nil {
		return nil, err
	}
	return &modDescriptionT{
		name:    name,
		version: *version2,
	}, nil
}

func (m *modDescriptionT) String() string {
	if m.version.Full == "" {
		return m.name
	} else {
		return fmt.Sprintf("%s==%s", m.name, m.version.Full)
	}
}

type modsFilesT struct {
	versions map[string][]modDescriptionT
	extra    map[string]bool
	missing  map[string]bool
}

func modsFiles() *modsFilesT {
	return &modsFilesT{
		versions: map[string][]modDescriptionT{},
		extra:    map[string]bool{},
		missing:  map[string]bool{},
	}
}

type modRelease struct {
	DownloadUrl string `json:"download_url"`
	SHA1        string
	FileName    string `json:"file_name"`
	Version     string
	InfoJson    struct {
		FactorioVersion string `json:"factorio_version"`
	} `json:"info_json"`
}

type modPortalResponse struct {
	Message  string
	Name     string
	Releases []modRelease
}

// ModCommandUsage ...
var ModCommandUsage = "Usage: $mod (add|remove|enable|disable) <modnames>"

// ModCommand returns the list of mods running on the server.
func ModCommand(s *discordgo.Session, args string) {
	argsList := strings.SplitN(args, " ", 2)
	if len(argsList) == 0 {
		support.SendFormat(s, ModCommandUsage)
		return
	}

	action := argsList[0]
	switch action {
	case "add", "remove", "enable", "disable":
		if len(argsList) < 2 {
			support.SendFormat(s, "Usage: $mod "+action+" <modname> [<modname>]+")
			return
		}
	default:
		support.SendFormat(s, ModCommandUsage)
		return
	}

	modnames, mismatched := support.QuoteSplit(strings.Join(argsList[1:], " "), "\"")
	if mismatched {
		support.Send(s, "Error: Mismatched quotes")
		return
	}
	var modDescriptions []modDescriptionT
	if action == "add" {
		for _, modname := range modnames {
			desc, err := modDescription(modname)
			if err != nil {
				support.Send(s, "Error parsing version: "+modname)
				return
			}
			modDescriptions = append(modDescriptions, *desc)
		}
	}

	modsListFile, err := ioutil.ReadFile(support.Config.ModListLocation)
	if err != nil {
		support.Send(s, "Sorry, there was an error reading your mod list")
		support.Panik(err, "there was an error reading mods list, did you specify it in the config.json file?")
		return
	}

	mods := &ModJSON{}
	err = json.Unmarshal(modsListFile, &mods)
	if err != nil {
		support.Send(s, "Sorry, there was an error reading your mod list")
		support.Panik(err, "there was an error reading mod list")
		return
	}

	var res string
	switch action {
	case "add":
		res = modsAdd(s, mods, &modDescriptions)
	case "remove":
		res = modsRemove(mods, modnames)
	case "enable":
		res = modsEnable(mods, modnames, true)
	case "disable":
		res = modsEnable(mods, modnames, false)
	}

	modsListFile, err = json.MarshalIndent(mods, "", "    ")
	if err != nil {
		support.Send(s, "Sorry, there was an error converting mod list to json")
		support.Panik(err, "there was an error converting mod list to json")
		return
	}
	err = ioutil.WriteFile(support.Config.ModListLocation, modsListFile, 0666)
	if err != nil {
		support.Send(s, "Sorry, there was an error saving mod list")
		support.Panik(err, "there was an error saving mod list")
		return
	}

	support.ChunkedMessageSend(s, res)
}

func modsAdd(s *discordgo.Session, mods *ModJSON, modDescriptions *[]modDescriptionT) string {
	modsList := make([]Mod, len(mods.Mods)+len(*modDescriptions))
	var toDownload []*modRelease

	files := matchModsWithFiles(&mods.Mods)

	end := len(mods.Mods)
	copy(modsList, mods.Mods)
	mods.Mods = modsList

	res := ""
	alreadyAdded := ""

	factorioVersion, err := support.FactorioVersion()
	if err != nil {
		return "Error checking factorio version"
	}
	factorioVersion = strings.Join(strings.Split(factorioVersion, ".")[:2], ".")

	userErrors := ""
	for _, desc := range *modDescriptions {
		if _, downloaded := files.versions[desc.name]; downloaded {
			alreadyAdded += "\n    " + desc.String()
			continue
		}
		release, userError, err := checkModPortal(&desc, factorioVersion)
		if err != nil {
			return "Some connection error occurred"
		}
		if userError != "" {
			userErrors += fmt.Sprintf("\n    %s: %s", desc.String(), userError)
			continue
		}
		toDownload = append(toDownload, release)
		added := false

		for i := 0; i < end; i++ {
			mod := mods.Mods[i]
			if strings.ToLower(mod.Name) == strings.ToLower(desc.name) {
				alreadyAdded += "\n    " + desc.String()
				added = true
				break
			}
			if strings.ToLower(mod.Name) > strings.ToLower(desc.name) {
				copy(mods.Mods[i+1:], mods.Mods[i:])
				mods.Mods[i] = Mod{
					Name:    desc.name,
					Enabled: true,
				}
				end++
				added = true
				res += "\n    " + desc.String()
				break
			}
		}
		if !added {
			res += "\n    " + desc.String()
			mods.Mods[end] = Mod{
				Name:    desc.name,
				Enabled: true,
				Version: desc.version.Full,
			}
			end++
		}
	}
	if len(*modDescriptions) == 1 {
		_, downloaded := files.versions[(*modDescriptions)[0].name]
		if alreadyAdded != "" && !downloaded {
			res = fmt.Sprintf("Mod \"%s\" is already added", (*modDescriptions)[0].String())
		} else if userErrors != "" {
			res = strings.TrimSpace(userErrors)
		} else {
			res = fmt.Sprintf("Added mod \"%s\"", (*modDescriptions)[0].String())
		}
	} else {
		res = "**Added mods:**" + res
		if alreadyAdded != "" {
			res += "\n**Already added:**" + alreadyAdded
		}
		if userErrors != "" {
			res += "\n**Errors:**" + userErrors
		}
	}
	mods.Mods = mods.Mods[:end]
	if !modDownloaderStarted {
		if support.Config.ModPortalToken == "" {
			res += "\n**No token to download mods**"
		} else if support.Config.Username == "" {
			res += "\n**No username to download mods**"
		} else {
			go modDownloader(s)
			for _, x := range toDownload {
				downloadQueue <- x
			}
		}
	} else {
		for _, x := range toDownload {
			downloadQueue <- x
		}
	}
	return res
}

func modsRemove(mods *ModJSON, modnames []string) string {
	removed := 0
	res := ""
	notFoundCount := 0
	notFound := ""
	removedFiles := ""
	errorRemovingFiles := false

	files := matchModsWithFiles(&mods.Mods)

	for _, modname := range modnames {
		found := false

		for i, mod := range mods.Mods[:len(mods.Mods)-removed] {
			if modname == mod.Name {
				found = true
				res += "\n    " + modname
				copy(mods.Mods[i:], mods.Mods[i+1:])
				removed++
				break
			}
		}
		if modFiles, ok := files.versions[modname]; ok {
			found = true
			if !errorRemovingFiles {
				for _, desc := range modFiles {
					err := os.Remove(desc.path)
					if err != nil {
						errorRemovingFiles = true
						removedFiles = "There was an error removing mod files. Try shutting down the server"
					}
					removedFiles = removedFiles + "\n    " + desc.String()
				}
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
		} else if removedFiles == "" {
			res = "Removed mod \"" + modnames[0] + "\""
		} else {
			res = "Removed " + strings.TrimSpace(removedFiles)
		}
	} else {
		res = fmt.Sprintf("**Removed %d mods (left: %d):**", removed, len(mods.Mods)) + res
		if removedFiles != "" {
			removedFiles = "\n**Files removed:**" + removedFiles
			res += removedFiles
		}
		if notFoundCount > 0 {
			notFound = fmt.Sprintf("\n**%d mods weren't found:**", notFoundCount) + notFound
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

func matchModsWithFiles(mods *[]Mod) *modsFilesT {
	res := modsFiles()
	for _, mod := range *mods {
		res.missing[mod.Name] = true
	}
	baseDir := path.Dir(support.Config.ModListLocation)
	files, err := ioutil.ReadDir(baseDir)
	if err != nil {
		support.Critical(err, "wtf")
	}

	for _, file := range files {
		re := support.ModFileRegexp.FindStringSubmatch(file.Name())
		if re == nil || re[1] == "" || re[2] == "" {
			continue
		}
		name := re[1]
		version, err := support.SemanticVersion(re[2])
		if err != nil {
			panic("wtf")
		}

		res.versions[name] = append(res.versions[name], modDescriptionT{
			name:    name,
			path:    path.Join(baseDir, file.Name()),
			version: *version,
		})

		found := false
		for _, mod := range *mods {
			if mod.Name == name {
				found = true
				delete(res.missing, mod.Name)
				break
			}
		}
		if !found {
			res.extra[name] = true
		}
	}
	return res
}

func checkModPortal(desc *modDescriptionT, factorioVersion string) (*modRelease, string, error) {
	resp, err := http.Get(fmt.Sprintf("https://mods.factorio.com/api/mods/%s", desc.name))
	if err != nil {
		return nil, "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, "", err
	}

	response := modPortalResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, "", err
	}
	if response.Message == "Mod not found" {
		return nil, "mod not found on the mod portal", nil
	}

	if desc.version.Full == "" { // no version specified
		for z := len(response.Releases) - 1; z >= 0; z-- {
			if response.Releases[z].InfoJson.FactorioVersion == factorioVersion {
				return &response.Releases[z], "", nil
			}
		}
		return nil, "no release for this factorio version", nil
	} else {
		for _, release := range response.Releases {
			if release.Version == desc.version.Full {
				if release.InfoJson.FactorioVersion == factorioVersion {
					return &release, "", nil
				} else {
					return nil, fmt.Sprintf(
						"this version of the mod (%s) is not for this factorio version (%s)",
						release.InfoJson.FactorioVersion,
						factorioVersion,
					), nil
				}
			}
		}
		return nil, "no such version", nil
	}
}

var downloadQueue = make(chan *modRelease, 100)
var modDownloaderStarted = false

func downloadProgressUpdater(s *discordgo.Session, wc *support.WriteCounter, modname string) {
	message := support.Send(s, fmt.Sprintf(support.Config.Messages.DownloadProgress, modname, wc.Percent()))
	time.Sleep(500 * time.Millisecond)
	for {
		if wc.Error {
			return
		}
		if wc.Progress >= wc.Total {
			break
		}
		message.Edit(s, fmt.Sprintf(support.Config.Messages.DownloadProgress, modname, wc.Percent()))
		time.Sleep(2 * time.Second)
	}
	message.Edit(s, fmt.Sprintf(support.Config.Messages.DownloadComplete, modname))
}

func modDownloader(s *discordgo.Session) {
	modDownloaderStarted = true
	baseDir := path.Dir(support.Config.ModListLocation)
	for {
		mod := <-downloadQueue

		file, err := os.OpenFile(
			path.Join(baseDir, mod.FileName),
			os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
			0664,
		)
		if err != nil {
			support.Panik(err, "Error opening "+mod.FileName+" for write")
			support.Send(s, mod.FileName+": error opening file for write")
		}

		url := fmt.Sprintf(
			"https://mods.factorio.com%s?username=%s&token=%s",
			mod.DownloadUrl,
			support.Config.Username,
			support.Config.ModPortalToken,
		)
		resp, err := http.Get(url)
		if err != nil {
			support.Panik(err, "Error downloading mod")
			support.Send(s, mod.FileName+": Error downloading mod")
			continue
		}
		if resp.ContentLength < 0 {
			if strings.Contains(resp.Request.URL.Path, "login") {
				support.Send(s, "Error logging in to download mods. Check username and mod portal token")
			} else {
				support.Panik(errors.New("content length error"), "Error downloading mod")
				support.Send(s, "Error downloading mod")
				continue
			}
		}

		counter := &support.WriteCounter{Total: uint64(resp.ContentLength)}
		go downloadProgressUpdater(s, counter, mod.FileName)
		_, err = io.Copy(io.MultiWriter(file, counter), resp.Body)
		file.Close()
		resp.Body.Close()
		if err != nil {
			counter.Error = true
			support.Panik(err, "Error downloading mod file")
			continue
		}
	}
}
