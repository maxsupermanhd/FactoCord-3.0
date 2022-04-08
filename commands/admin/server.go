package admin

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/maxsupermanhd/FactoCord-3.0/support"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var ServerCommandDoc = support.Command{
	Name: "server",
	Desc: "Manage factorio server",
	Usage: "/server\n" +
		"/server [stop|start|restart|update <version>?]",
	Doc: "command manages factorio server.\n" +
		"`/server` shows current server status. Anyone can execute it.`",
	Admin:   false,
	Command: nil,
	Subcommands: []support.Command{
		{
			Name: "status",
			Desc: `Show current server status`,
			Doc:  `command shows current server status`,
			Command: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				if support.Factorio.IsRunning() {
					support.Respond(s, i, "Factorio server is **running**")
				} else {
					support.Respond(s, i, "Factorio server is **stopped**")
				}
			},
		},
		{
			Name:  "stop",
			Desc:  `Stop the server`,
			Doc:   `command stops the server`,
			Admin: true,
			Command: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				support.Factorio.Stop(s)
				support.Respond(s, i, "Stopping server")
			},
		},
		{
			Name:  "start",
			Desc:  `Start the server`,
			Doc:   `command starts the server`,
			Admin: true,
			Command: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				support.Factorio.Start(s)
				support.Respond(s, i, "Starting server")
			},
		},
		{
			Name:  "restart",
			Desc:  `Restart the server`,
			Doc:   `command restarts the server`,
			Admin: true,
			Command: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				support.Factorio.Stop(s)
				support.Factorio.Start(s)
				support.Respond(s, i, "Restarting server")
			},
		},
		{
			Name: "update",
			Desc: `Update to server to the newest version or to the specified version`,
			Doc:  `command updates to server to the newest version or to the specified version`,
			Usage: "/server update\n" +
				"/server update <version>",
			Admin: true,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "version",
					Description: "Factorio version to update to",
					Required:    false,
				},
			},
			Command: serverUpdate,
		},
	},
}

func serverUpdate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	version := ""
	options := i.ApplicationCommandData().Options[0].Options
	if len(options) == 1 {
		version = i.ApplicationCommandData().Options[0].Options[0].StringValue()
	}
	if support.Factorio.IsRunning() {
		support.Send(s, "You should `/server stop` first")
		return
	}
	//username := support.Config.Username
	//token := support.Config.ModPortalToken
	//if username == "" {
	//	support.Send(s, "Username is required for update")
	//	return
	//}
	//if token == "" {
	//	support.Send(s, "Token is required for update")
	//	return
	//}
	factorioVersion, err := support.FactorioVersion()
	if err != nil {
		support.Panik(err, "... checking factorio version")
		support.Respond(s, i, "Error checking factorio version")
		return
	}
	support.RespondDefer(s, i, "Getting versions...")

	if version == "" {
		version, err = getLatestVersion()
		if err != nil {
			support.Panik(err, "Error getting latest version information")
			support.ResponseEdit(s, i, "Error getting latest version information")
			return
		}
		if version == factorioVersion {
			support.ResponseEdit(s, i, "The server is already updated to the latest version")
			return
		}
	} else if version == factorioVersion {
		support.ResponseEdit(s, i, "The server is already updated to that version")
		return
	}

	resp, err := http.Get(fmt.Sprintf("https://updater.factorio.com/get-download/%s/headless/linux64", version))
	if err != nil {
		support.Panik(err, "Connection error downloading factorio")
		support.ResponseEdit(s, i, "Some connection error occurred")
		return
	}
	if resp.StatusCode == 404 {
		support.ResponseEdit(s, i, fmt.Sprintf("Version %s not found\n"+
			"Refer to <https://factorio.com/download/archive> to see available versions", version))
		return
	}
	if resp.ContentLength <= 0 {
		support.ResponseEdit(s, i, "Error with content-length")
		return
	}

	urlPath := strings.Split(resp.Request.URL.Path, "/")
	filename := urlPath[len(urlPath)-1]
	if filename == "" {
		support.ResponseEdit(s, i, "Error with filename")
		support.Panik(fmt.Errorf("uncorrect url path: %s", resp.Request.URL.Path), "Error with filename")
		return
	}
	path := "/tmp/" + filename

	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664)
	if err != nil {
		support.Panik(err, "Error opening "+path+" for write")
		support.ResponseEdit(s, i, path+": error opening file for write")
		return
	}

	support.ResponseEdit(s, i, support.FormatNamed(support.Config.Messages.DownloadStart, "file", filename))
	counter := &support.WriteCounter{Total: uint64(resp.ContentLength)}
	progress := support.ProgressUpdate{
		WriteCounter: counter,
		Interaction:  i,
		Progress:     support.FormatNamed(support.Config.Messages.DownloadProgress, "file", filename),
		Finished:     support.FormatNamed(support.Config.Messages.Unpacking, "file", filename),
	}
	go support.DownloadProgressUpdater(s, &progress)

	_, err = io.Copy(io.MultiWriter(file, counter), resp.Body)
	resp.Body.Close()
	file.Close()
	if err != nil {
		counter.Error = true
		support.ResponseEdit(s, i, ":interrobang: Error downloading "+filename)
		support.Panik(err, "Error downloading file")
		return
	}

	dir, err := filepath.Abs(support.Config.Executable)
	if err != nil {
		support.Panik(err, "Error getting absolute path of executable")
		support.ResponseEdit(s, i, "Error getting absolute path of executable")
		return
	}
	dir = filepath.Dir(dir) // x64
	dir = filepath.Dir(dir) // bin
	dir = filepath.Dir(dir) // factorio
	cmd := exec.Command("tar", "-C", dir, "--strip-components=1", "-xf", path)
	err = cmd.Run()
	if err != nil {
		support.Panik(err, "Error running tar to unpack the archive")
		support.ResponseEdit(s, i, "Error running tar to unpack the archive")
		return
	}

	support.ResponseEdit(s, i, support.FormatNamed(support.Config.Messages.UnpackingComplete, "version", version))
	_ = os.Remove(path)
}

type latestVersions struct {
	Stable, Experimental struct {
		Alpha, Demo, Headless string
	}
}

func getLatestVersion() (string, error) {
	resp, err := http.Get("https://factorio.com/api/latest-releases")
	if err != nil {
		return "", err
	}
	var versions latestVersions
	err = json.NewDecoder(resp.Body).Decode(&versions)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	return versions.Experimental.Headless, nil
}
