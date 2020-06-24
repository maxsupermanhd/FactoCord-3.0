package support

import (
	"fmt"
	"github.com/yosuke-furukawa/json5/encoding/json5"
	"io/ioutil"
	"os"
)

// Config is a config interface.
var Config configT

type configT struct {
	Executable       string   `json:"executable"`
	LaunchParameters []string `json:"launch_parameters"`

	DiscordToken            string `json:"discord_token"`
	GameName                string `json:"game_name"`
	FactorioChannelID       string `json:"factorio_channel_id"`
	Prefix                  string `json:"prefix"`
	HaveServerEssentials    bool   `json:"have_server_essentials"`
	IngameDiscordUserColors bool   `json:"ingame_discord_user_colors"`

	EnableConsoleChannel  bool   `json:"enable_console_channel"`
	FactorioConsoleChatID string `json:"factorio_console_chat_id"`

	AdminIDs        []string          `json:"admin_ids"`
	CommandRoles    map[string]string `json:"command_roles"`
	ModListLocation string            `json:"mod_list_location"`

	Messages struct {
		BotStart    string `json:"bot_start"`
		BotStop     string `json:"bot_stop"`
		ServerStart string `json:"server_start"`
		ServerStop  string `json:"server_stop"`
		ServerFail  string `json:"server_fail"`
		ServerSave  string `json:"server_save"`
		PlayerJoin  string `json:"player_join"`
		PlayerLeave string `json:"player_leave"`
	} `json:"messages"`
}

func (conf *configT) Load() {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		fmt.Println("Environment file not found, cannot continue!")
		Exit(1)
	}
	if !FileExists("./config.json") {
		fmt.Println("Error: config.json not found.")
		fmt.Println("Make sure that you copied 'config-example.json' and current working directory is correct")
		Exit(7)
	}
	contents, err := ioutil.ReadFile("./config.json")
	Critical(err, "... when reading config.json")
	err = json5.Unmarshal(contents, &conf)
	if err != nil {
		fmt.Println("Note that json5 may have several bugs, such as comment before ] or }")
		Critical(err, "... when parsing config.json")
	}
}
