package admin

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

// Mod is a struct containing info about a mod.
type Mod struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Version string `json:"version,omitempty"`
}

func (m *Mod) Description() *modDescriptionT {
	version, err := support.SemanticVersion(m.Version)
	if err != nil {
		panic(err)
	}
	return &modDescriptionT{
		name:    m.Name,
		path:    "",
		version: *version,
	}
}

// ModJSON is struct containing a slice of Mod.
type ModJSON struct {
	Mods []Mod `json:"mods"`
}

func (m *ModJSON) sortedInsert(newMod *Mod) bool {
	for i := 0; i < len(m.Mods); i++ {
		mod := m.Mods[i]
		if strings.ToLower(mod.Name) == strings.ToLower(newMod.Name) {
			return false
		}
		if strings.ToLower(mod.Name) > strings.ToLower(newMod.Name) {
			m.Mods = append(m.Mods, Mod{})
			copy(m.Mods[i+1:], m.Mods[i:])
			m.Mods[i] = *newMod
			return true
		}
	}
	m.Mods = append(m.Mods, *newMod)
	return true
}

func (m *ModJSON) removeMod(modname string) (removed bool) {
	for i, mod := range m.Mods {
		if modname == mod.Name {
			copy(m.Mods[i:], m.Mods[i+1:])
			m.Mods = m.Mods[:len(m.Mods)-1]
			return true
		}
	}
	return false
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

func (m *modDescriptionT) ModEntry() *Mod {
	return &Mod{
		Name:    m.name,
		Enabled: true,
		Version: m.version.Full,
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
var ModCommandUsage = "Usage: $mod (add|remove|enable|disable) <modnames>+ | update <modnames>*"

// ModCommand returns the list of mods running on the server.
func ModCommand(s *discordgo.Session, args string) {
	argsList := strings.SplitN(args, " ", 2)
	if len(argsList) == 0 {
		support.SendFormat(s, ModCommandUsage)
		return
	}

	action := argsList[0]
	switch action {
	case "update":
		//
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
	if action == "add" || action == "update" {
		for _, modname := range modnames {
			desc, err := modDescription(modname)
			if err != nil {
				support.Send(s, "Error parsing version: "+modname)
				return
			}
			modDescriptions = append(modDescriptions, *desc)
		}
		var t []interface{}
		for _, x := range modDescriptions {
			t = append(t, x) // some golang shit
		}
		if support.AnyTwo(t, func(desc interface{}, desc2 interface{}) bool {
			return desc.(modDescriptionT).name == desc2.(modDescriptionT).name
		}) {
			support.Send(s, "Who am I supposed to add a single mod twice?")
			return
		}
	} else if !support.IsUnique(modnames) {
		support.Send(s, "Who am I supposed to change a single mod twice?")
		return
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
	case "update":
		res = modsUpdate(s, mods, &modDescriptions)
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
	var toDownload []*modRelease

	files := matchModsWithFiles(&mods.Mods)

	addedMods := support.DefaultTextList("**Added mods:**")
	alreadyAdded := support.DefaultTextList("\n**Already added:**")
	userErrors := support.DefaultTextList("\n**Errors:**")

	factorioVersion, err := factorioVersion()
	if err != nil {
		return "Error checking factorio version"
	}

	for _, desc := range *modDescriptions {
		if _, downloaded := files.versions[desc.name]; downloaded {
			alreadyAdded.Append(desc.String())
			continue
		}
		release, userError, err := checkModPortal(&desc, factorioVersion)
		if err != nil {
			return "Some connection error occurred"
		}
		if userError != "" {
			userErrors.Append(fmt.Sprintf("%s: %s", desc.String(), userError))
			continue
		}

		toDownload = append(toDownload, release)
		inserted := mods.sortedInsert(desc.ModEntry())
		if inserted {
			addedMods.Append(desc.String())
		} else {
			alreadyAdded.Append(desc.String())
			if desc.version.Full != "" {
				alreadyAdded.AddToLast(support.FormatUsage(" - to update a mod use `$mod update` command"))
			}
		}
	}
	res := ""
	if len(*modDescriptions) == 1 {
		_, downloaded := files.versions[(*modDescriptions)[0].name]
		if alreadyAdded.NotEmpty() && !downloaded {
			res = fmt.Sprintf("Mod \"%s\" is already added", (*modDescriptions)[0].String())
		} else if userErrors.NotEmpty() {
			res = strings.TrimSpace(userErrors.List[0])
		} else {
			res = fmt.Sprintf("Added mod \"%s\"", (*modDescriptions)[0].String())
		}
	} else {
		res = addedMods.Render()
		res += alreadyAdded.RenderNotEmpty()
		res += userErrors.RenderNotEmpty()
	}

	if support.Config.ModPortalToken == "" {
		res += "\n**No token to download mods**"
	} else if support.Config.Username == "" {
		res += "\n**No username to download mods**"
	} else {
		if !modDownloaderStarted {
			go modDownloader(s)
		}
		for _, x := range toDownload {
			downloadQueue <- x
		}
	}
	return res
}

func factorioVersion() (string, error) {
	factorioVersion, err := support.FactorioVersion()
	if err != nil {
		return "", err
	}
	factorioVersion = strings.Join(strings.Split(factorioVersion, ".")[:2], ".")
	return factorioVersion, nil
}

func modsUpdate(s *discordgo.Session, mods *ModJSON, modDescriptions *[]modDescriptionT) string {
	if support.Config.ModPortalToken == "" {
		return "**No token to download mods**"
	} else if support.Config.Username == "" {
		return "**No username to download mods**"
	}

	updatedMods := support.DefaultTextList("**Updating mods:**")
	alreadyUpdated := support.DefaultTextList("\n**Already updated:**")
	userErrors := support.DefaultTextList("\n**Errors:**")

	var toDownload []*modRelease

	files := matchModsWithFiles(&mods.Mods)

	factorioVersion, err := factorioVersion()
	if err != nil {
		return "Error checking factorio version"
	}

	updateAll := true
	if len(*modDescriptions) == 0 {
		updateAll = false
		*modDescriptions = nil
		for _, mod := range mods.Mods {
			if mod.Name != "base" {
				*modDescriptions = append(*modDescriptions, modDescriptionT{name: mod.Name})
			}
		}
	}

	for _, desc := range *modDescriptions {
		release, userError, err := checkModPortal(&desc, factorioVersion)
		if err != nil {
			return "Some connection error occurred"
		}
		if userError != "" {
			userErrors.Append(fmt.Sprintf("%s: %s", desc.String(), userError))
			continue
		}

		versions := files.versions[desc.name]
		var versionsVersions []support.SemanticVersionT
		var versionsStrings []string
		downloaded := false
		for _, version := range versions {
			versionsStrings = append(versionsStrings, version.version.Full)
			versionsVersions = append(versionsVersions, version.version)
			if version.version.Full == release.Version {
				downloaded = true
			}
		}
		if downloaded {
			alreadyUpdated.Append(desc.String())
			continue
		}
		releaseVersion := support.SemanticVersionPanic(release.Version)
		toDownload = append(toDownload, release)
		updatedMods.Append(fmt.Sprintf(
			"**%s** %s **%s %s**",
			desc.name,
			strings.Join(versionsStrings, ", "),
			versionsArrow(versionsVersions, releaseVersion),
			release.Version,
		))
		_, err = removeModFiles(files, desc.name)
		if err != nil {
			updatedMods.AddToLast(": error removing files")
		}
	}
	if !modDownloaderStarted {
		go modDownloader(s)
	}
	for _, x := range toDownload {
		downloadQueue <- x
	}
	if updateAll {
		return updatedMods.Render() + alreadyUpdated.RenderNotEmpty() + userErrors.RenderNotEmpty()
	} else {
		return updatedMods.Render() + userErrors.RenderNotEmpty()
	}
}

func modsRemove(mods *ModJSON, modnames []string) string {
	removedMods := support.DefaultTextList("**Removed %d mods (left: %d):**")
	notFound := support.DefaultTextList("\n**%d mods weren't found:**")
	removedFiles := support.DefaultTextList("\n**Files removed:**")

	files := matchModsWithFiles(&mods.Mods)

	for _, modname := range modnames {
		found := mods.removeMod(modname)
		if found {
			removedMods.Append(modname)
		}

		if removedFiles.Error == "" {
			filesFound, err := removeModFiles(files, modname)
			if err != nil {
				removedFiles.Error = "\nThere was an error removing mod files. Try shutting down the server"
				continue
			}
			found = found || len(filesFound) > 0
			for _, desc := range filesFound {
				removedFiles.Append(desc.String())
			}
		}
		if !found {
			notFound.Append(modname)
		}
	}
	if len(modnames) == 1 {
		if notFound.NotEmpty() {
			return "Mod \"" + modnames[0] + "\" not found"
		} else if removedFiles.NotEmpty() {
			if removedFiles.Error != "" {
				return removedFiles.Error
			}
			return "Removed " + removedFiles.List[0]
		} else {
			return "Removed mod \"" + modnames[0] + "\""
		}
	} else {
		removedMods.Heading = fmt.Sprintf(removedMods.Heading, removedMods.Len(), len(mods.Mods))
		notFound.FormatHeaderWithLength()
		return removedMods.Render() + removedFiles.RenderNotEmpty() + notFound.RenderNotEmpty()
	}
}

func modsEnable(mods *ModJSON, modnames []string, setEnabled bool) string {
	toggled := support.DefaultTextList("")
	notFound := support.DefaultTextList("\n**Not Found %d mods:**")

	for _, modname := range modnames {
		found := false
		for i, mod := range mods.Mods {
			if mod.Name == modname {
				mods.Mods[i].Enabled = setEnabled
				found = true
				toggled.Append(modname)
			}
		}
		if !found {
			notFound.Append(modname)
		}
	}

	action := "Disabled"
	if setEnabled {
		action = "Enabled"
	}
	if len(modnames) == 1 {
		if notFound.NotEmpty() {
			return "Mod \"" + modnames[0] + "\" not found"
		} else {
			return action + " mod \"" + modnames[0] + "\""
		}
	} else {
		toggled.Heading = "**" + action + " %d mods:**"
		toggled.FormatHeaderWithLength()
		notFound.FormatHeaderWithLength()
		return toggled.Render() + notFound.RenderNotEmpty()
	}
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

func removeModFiles(files *modsFilesT, modname string) (found []modDescriptionT, err error) {
	modFiles, ok := files.versions[modname]
	if !ok {
		return nil, nil
	}
	for _, desc := range modFiles {
		err := os.Remove(desc.path)
		if err != nil {
			return modFiles, err
		}
	}
	return modFiles, nil
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

func modDownloader(s *discordgo.Session) {
	modDownloaderStarted = true
	baseDir := path.Dir(support.Config.ModListLocation)
	for {
		mod := <-downloadQueue

		file, err := os.OpenFile(
			path.Join(baseDir, mod.FileName),
			os.O_CREATE|os.O_TRUNC|os.O_RDWR,
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

		message := support.Send(s, support.FormatNamed(support.Config.Messages.DownloadStart, "file", mod.FileName))
		counter := &support.WriteCounter{Total: uint64(resp.ContentLength)}
		progress := support.ProgressUpdate{
			WriteCounter: counter,
			Message:      message,
			Progress:     support.FormatNamed(support.Config.Messages.DownloadProgress, "file", mod.FileName),
			Finished:     support.FormatNamed(support.Config.Messages.DownloadComplete, "file", mod.FileName),
		}
		go support.DownloadProgressUpdater(s, &progress)

		_, err = io.Copy(io.MultiWriter(file, counter), resp.Body)
		resp.Body.Close()
		if err != nil {
			counter.Error = true
			support.Panik(err, "Error downloading mod file")
			continue
		}

		if mod.SHA1 != "" {
			_, err = file.Seek(0, 0) // to the start
			if err != nil {
				panic(err)
			}

			hash, err := fileHash(file)
			if err != nil {
				support.Panik(err, "... calculating sha1")
				continue
			}
			if mod.SHA1 != hash {
				counter.Error = true
				message.Edit(s, fmt.Sprintf(":interrobang: %s is downloaded but hashsum is invalid", mod.FileName))
			}
		}
		file.Close()
	}
}

func fileHash(file io.Reader) (string, error) {
	hash := sha1.New()
	_, err := io.Copy(hash, file)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func versionsArrow(v1 []support.SemanticVersionT, v2 *support.SemanticVersionT) string {
	if len(v1) == 1 {
		if v2.NewerThan(&v1[0]) {
			return "⭧"
		} else {
			return "⭨"
		}
	} else {
		return "⭢"
	}
}
