package admin

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/maxsupermanhd/FactoCord-3.0/v3/support"
)

var ServerCommandDoc = support.CommandDoc{
	Name: "server",
	Usage: "$server\n" +
		"$server [stop|start|restart|update <version>?]",
	Doc: "command manages factorio server.\n" +
		"`$server` shows current server status. Anyone can execute it.`",
	Subcommands: []support.CommandDoc{
		{Name: "stop", Doc: `command stops the server`},
		{Name: "start", Doc: `command starts the server`},
		{Name: "restart", Doc: `command restarts the server`},
		{
			Name: "update",
			Doc:  `command updates to server to the newest version or to the specified version`,
			Usage: "$server update\n" +
				"$server update <version>",
		},
	},
}

func ServerCommandAdminPermission(args string) bool {
	return strings.TrimSpace(args) != ""
}

func ServerCommand(s *discordgo.Session, args string) {
	action, arg := support.SplitDivide(args, " ")
	switch action {
	case "":
		if support.Factorio.IsRunning() {
			support.Send(s, "Factorio server is **running**")
		} else {
			support.Send(s, "Factorio server is **stopped**")
		}
	case "stop":
		support.Factorio.Stop(s)
	case "start":
		support.Factorio.Start(s)
	case "restart":
		support.Factorio.Stop(s)
		support.Factorio.Start(s)
	case "update":
		serverUpdate(s, arg)
	default:
		support.SendFormat(s, "Usage: "+ServerCommandDoc.Usage)
	}
}

func serverUpdate(s *discordgo.Session, version string) {
	if support.Factorio.IsRunning() {
		support.Send(s, "You should stop the server first")
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
		support.Send(s, "Error checking factorio version")
		return
	}

	if version == "" {
		version, err = getLatestVersion()
		if err != nil {
			support.Panik(err, "Error getting latest version information")
			support.Send(s, "Error getting latest version information")
			return
		}
		if version == factorioVersion {
			support.Send(s, "The server is already updated to the latest version")
			return
		}
	} else if version == factorioVersion {
		support.Send(s, "The server is already updated to that version")
		return
	}

	resp, err := http.Get(fmt.Sprintf("https://updater.factorio.com/get-download/%s/headless/linux64", version))
	if err != nil {
		support.Panik(err, "Connection error downloading factorio")
		support.Send(s, "Some connection error occurred")
		return
	}
	if resp.StatusCode == 404 {
		support.Send(s, fmt.Sprintf("Version %s not found\n"+
			"Refer to <https://factorio.com/download/archive> to see available versions", version))
		return
	}
	if resp.ContentLength <= 0 {
		support.Send(s, "Error with content-length")
		return
	}
	if resp.Header.Get("Content-Disposition") == "" {
		support.Send(s, "Error with content-disposition")
		return
	}
	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	if err != nil {
		support.Send(s, "Error with content-disposition")
		return
	}
	filename, ok := params["filename"]
	if !ok {
		support.Send(s, "Error with content-disposition")
		return
	}
	path := "/tmp/" + filename

	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664)
	if err != nil {
		support.Panik(err, "Error opening "+path+" for write")
		support.Send(s, path+": error opening file for write")
		return
	}

	message := support.Send(s, support.FormatNamed(support.Config.Messages.DownloadStart, "file", filename))
	counter := &support.WriteCounter{Total: uint64(resp.ContentLength)}
	progress := support.ProgressUpdate{
		WriteCounter: counter,
		Message:      message,
		Progress:     support.FormatNamed(support.Config.Messages.DownloadProgress, "file", filename),
		Finished:     support.FormatNamed(support.Config.Messages.Unpacking, "file", filename),
	}
	go support.DownloadProgressUpdater(s, &progress)

	_, err = io.Copy(io.MultiWriter(file, counter), resp.Body)
	resp.Body.Close()
	file.Close()
	if err != nil {
		counter.Error = true
		message.Edit(s, ":interrobang: Error downloading "+filename)
		support.Panik(err, "Error downloading file")
		return
	}

	dir, err := filepath.Abs(support.Config.Executable)
	if err != nil {
		support.Panik(err, "Error getting absolute path of executable")
		support.Send(s, "Error getting absolute path of executable")
		return
	}
	dir = filepath.Dir(dir) // x64
	dir = filepath.Dir(dir) // bin
	dir = filepath.Dir(dir) // factorio
	cmd := exec.Command("tar", "-C", dir, "--strip-components=1", "-xf", path)
	err = cmd.Run()
	if err != nil {
		support.Panik(err, "Error running tar to unpack the archive")
		support.Send(s, "Error running tar to unpack the archive")
		return
	}

	message.Edit(s, support.FormatNamed(support.Config.Messages.UnpackingComplete, "version", version))
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
