package utils

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/maxsupermanhd/FactoCord-3.0/support"
	"io/ioutil"
	"net/http"
	"strings"
)

var InfoDoc = support.Command{
	Name:    "info",
	Desc:    "Get server info",
	Doc:     `shows info about the server as from the factorio lobby`,
	Command: GameInfo,
}

type gameInfo struct {
	Message            string `json:"message"`
	ApplicationVersion struct {
		BuildMode    string `json:"build_mode"`
		BuildVersion int    `json:"build_version"`
		GameVersion  string `json:"game_version"`
		Platform     string `json:"platform"`
	} `json:"application_version"`
	Description     string  `json:"description"`
	GameID          int     `json:"game_id"`
	GameTimeElapsed int     `json:"game_time_elapsed"`
	HasPassword     bool    `json:"has_password"`
	HostAddress     string  `json:"host_address"`
	LastHeartbeat   float64 `json:"last_heartbeat"`
	MaxPlayers      int     `json:"max_players"`
	Mods            []struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"mods"`
	ModsCrc  int      `json:"mods_crc"`
	Name     string   `json:"name"`
	Players  []string `json:"players"`
	ServerID string   `json:"server_id"`
	Tags     []string `json:"tags"`
}

func GameInfo(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !support.Factorio.IsRunning() {
		support.Respond(s, i, "The server is not running")
		return
	}
	if support.Factorio.GameID == "" {
		support.Respond(s, i, "The server did not register a game on the factorio server")
		return
	}
	support.RespondDefer(s, i, "Fetching...")

	resp, err := http.Get("https://multiplayer.factorio.com/get-game-details/" + support.Factorio.GameID)
	if err != nil {
		support.Panik(err, "Connection error to /get-game-details")
		support.ResponseEdit(s, i, "Some connection error occurred")
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		support.Panik(err, "Error reading /get-game-details")
		support.ResponseEdit(s, i, "Some connection error occurred")
		return
	}

	info := gameInfo{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		support.Panik(err, "Error unmarshalling /get-game-details")
		support.ResponseEdit(s, i, "Some json error occurred")
		return
	}
	if info.Message != "" {
		support.ResponseEdit(s, i, "The server reports: "+info.Message)
		return
	}
	embed := &discordgo.MessageEmbed{
		Type:        "rich",
		Color:       0x6289FF,
		Title:       info.Name,
		Description: info.Description,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Version",
				Value: info.ApplicationVersion.GameVersion,
			},
		},
	}
	if len(info.Tags) > 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "Tags",
			Value: strings.Join(info.Tags, "\n"),
		})
	}
	online := getOnline(&info)
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:  online.Heading,
		Value: online.RenderWithoutHeading(),
	})
	support.ResponseEditCompex(s, i, embed)
}
