package utils

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"

	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

var OnlineDoc = support.Command{
	Name:    "online",
	Desc:    "Get players online",
	Doc:     `shows players online (and max number of players if set)`,
	Command: GameOnline,
}

func getOnline(info *gameInfo) *support.TextListT {
	if len(info.Players) == 0 {
		return &support.TextListT{
			Heading: "**No one is online**",
			None:    "",
		}
	}
	maxPlayers := ""
	if info.MaxPlayers != 0 {
		maxPlayers = fmt.Sprintf("/%d", info.MaxPlayers)
	}
	online := support.DefaultTextList(
		fmt.Sprintf("**%d%s player%s online:**", len(info.Players), maxPlayers, support.PluralS(len(info.Players))),
	)
	for _, player := range info.Players {
		online.Append(player)
	}
	return &online
}

func GameOnline(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	support.ResponseEdit(s, i, getOnline(&info).Render())
}
