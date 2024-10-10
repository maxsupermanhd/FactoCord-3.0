package support

import (
	"fmt"
	"os"

	"github.com/titanous/json5"
)

var FactoCordVersion string

var ConfigPath = "./config.json"

var GuildID string

// Config is a config interface.
var Config configT

type configT struct {
	Executable       string   `json:"executable"`
	LaunchParameters []string `json:"launch_parameters"`
	Autolaunch       bool     `json:"autolaunch"`

	DiscordToken            string `json:"discord_token"`
	GameName                string `json:"game_name"`
	FactorioChannelID       string `json:"factorio_channel_id"`
	Prefix                  string `json:"prefix"`
	HaveServerEssentials    bool   `json:"have_server_essentials"`
	IngameDiscordUserColors bool   `json:"ingame_discord_user_colors"`

	AllowPingingEveryone bool `json:"allow_pinging_everyone"`

	EnableConsoleChannel  bool   `json:"enable_console_channel"`
	FactorioConsoleChatID string `json:"factorio_console_chat_id"`

	AdminIDs     []string          `json:"admin_ids"`
	CommandRoles map[string]string `json:"command_roles"`

	ModListLocation string `json:"mod_list_location"`
	Username        string `json:"username"`
	ModPortalToken  string `json:"mod_portal_token"`

	Messages struct {
		BotStartLaunch    string `json:"bot_start"`
		BotStartOnly      string `json:"bot_start_only"`
		BotStop           string `json:"bot_stop"`
		ServerStart       string `json:"server_start"`
		ServerStop        string `json:"server_stop"`
		ServerFail        string `json:"server_fail"`
		ServerSave        string `json:"server_save"`
		PlayerJoin        string `json:"player_join"`
		PlayerLeave       string `json:"player_leave"`
		DownloadStart     string `json:"download_start"`
		DownloadProgress  string `json:"download_progress"`
		DownloadComplete  string `json:"download_complete"`
		Unpacking         string `json:"unpacking"`
		UnpackingComplete string `json:"unpacking_complete"`
	} `json:"messages"`
}

func (conf *configT) MustLoad() {
	if !FileExists(ConfigPath) {
		fmt.Println("Error: config.json not found.")
		fmt.Println("Make sure that you copied 'config-example.json' and current working directory is correct")
		Exit(7)
	}
	contents, err := os.ReadFile(ConfigPath)
	Critical(err, "... when reading config.json")

	conf.defaults()
	err = json5.Unmarshal(contents, &conf)
	if err != nil {
		Critical(err, "... when parsing config.json")
	}
}

func (conf *configT) Load() error {
	if !FileExists(ConfigPath) {
		return fmt.Errorf("config.json not found")
	}
	contents, err := os.ReadFile(ConfigPath)
	if err != nil {
		return fmt.Errorf("error reading config.json: %s", err)
	}

	test := configT{}
	err = json5.Unmarshal(contents, &test)
	if err != nil {
		return fmt.Errorf("error parsing config.json: %s", err)
	}

	conf.defaults()
	err = json5.Unmarshal(contents, &conf)
	Critical(err, "wtf?? error parsing config.json 2nd time")
	return nil
}

func (conf *configT) defaults() {
	conf.Autolaunch = true
	conf.GameName = "Factorio"
	conf.Prefix = "$"
	// conf.HaveServerEssentials = false
	// conf.IngameDiscordUserColors = false
	conf.Messages.BotStartLaunch = "**:white_check_mark: Bot started! Launching server...**"
	conf.Messages.BotStartOnly = "**:white_check_mark: Bot started! Autolaunch disabled.**"
	conf.Messages.BotStop = ":v:"
	conf.Messages.ServerStart = "**:white_check_mark: The server has started!**"
	conf.Messages.ServerStop = "**:octagonal_sign: The server has stopped!**"
	conf.Messages.ServerFail = "**:skull: The server has crashed!**"
	conf.Messages.ServerSave = "**:floppy_disk: Game saved!**"
	conf.Messages.PlayerJoin = "**:arrow_up: {username}**"
	conf.Messages.PlayerLeave = "**:arrow_down: {username}s**"
	conf.Messages.DownloadStart = ":arrow_down: Downloading {file}..."
	conf.Messages.DownloadProgress = ":arrow_down: Downloading {file}: {percent}%"
	conf.Messages.DownloadComplete = ":white_check_mark: Downloaded {file}"
	conf.Messages.Unpacking = ":pinching_hand: Unpacking {file}..."
	conf.Messages.UnpackingComplete = ":ok_hand: Server updated to {version}"
}
