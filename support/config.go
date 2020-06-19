package support

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config is a config interface.
var Config config

var ErrEnvVarEmpty = errors.New("getenv: environment variable empty")

type commandRolesType map[string]string

type config struct {
	DiscordToken          string
	FactorioChannelID     string
	Executable            string
	LaunchParameters      []string
	AdminIDs              []string
	CommandRoles          commandRolesType
	Prefix                string
	ModListLocation       string
	GameName              string
	PassConsoleChat       bool
	EnableConsoleChannel  bool
	FactorioConsoleChatID string
	HaveServerEssentials  bool
	BotStart              string
	SendBotStart          bool
	BotStop               string
	ServerStart           string
	ServerStop            string
	ServerFail            string
	ServerSave            string
	PlayerJoin            string
	PlayerLeave           string
}

func getenvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, ErrEnvVarEmpty
	}
	return v, nil
}

func getenvBool(key string) bool {
	s, err := getenvStr(key)
	if err != nil {
		return false
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return v
}

func getRolesMap(key string) commandRolesType {
	commandRoles := make(commandRolesType)
	s := os.Getenv(key)
	for _, roleCommands := range strings.Split(s, ";") {
		roleCommands = strings.TrimSpace(roleCommands)
		if roleCommands == "" {
			continue
		}
		roleCommandsSplit := strings.SplitN(roleCommands, ":", 2)
		if len(roleCommandsSplit) < 2 {
			fmt.Println(".env CommandRoles: Error parsing role with commands")
			Exit(1)
		}
		role := strings.TrimSpace(roleCommandsSplit[0])
		commands := strings.Split(strings.TrimSpace(roleCommandsSplit[1]), ",")
		if num, err := strconv.Atoi(role); err != nil || num <= 0 {
			fmt.Println(".env CommandRoles: " + key + " is incorrect: role is not an integer")
			Exit(1)
		}

		for _, command := range commands {
			if command == "" {
				continue
			}
			command = strings.ToLower(command)
			if _, exists := commandRoles[command]; exists {
				fmt.Println(".env CommandRoles: Command \"" + command + "\" is assigned multiple roles")
				Exit(1)
			}
			commandRoles[command] = role
		}
	}
	return commandRoles
}

func (conf *config) LoadEnv() {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		fmt.Println("Environment file not found, cannot continue!")
		Exit(1)
	}
	Config = config{
		DiscordToken:          os.Getenv("DiscordToken"),
		FactorioChannelID:     os.Getenv("FactorioChannelID"),
		LaunchParameters:      strings.Split(os.Getenv("LaunchParameters"), " "),
		Executable:            os.Getenv("Executable"),
		AdminIDs:              strings.Split(os.Getenv("AdminIDs"), ","),
		CommandRoles:          getRolesMap("CommandRoles"),
		Prefix:                os.Getenv("Prefix"),
		ModListLocation:       os.Getenv("ModListLocation"),
		GameName:              os.Getenv("GameName"),
		PassConsoleChat:       getenvBool("PassConsoleChat"),
		EnableConsoleChannel:  getenvBool("EnableConsoleChannel"),
		FactorioConsoleChatID: os.Getenv("FactorioConsoleChatID"),
		HaveServerEssentials:  getenvBool("HaveServerEssentials"),
		BotStart:              os.Getenv("BotStart"),
		SendBotStart:          getenvBool("SendBotStart"),
		BotStop:               os.Getenv("BotStop"),
		ServerStart:           os.Getenv("ServerStart"),
		ServerStop:            os.Getenv("ServerStop"),
		ServerFail:            os.Getenv("ServerFail"),
		ServerSave:            os.Getenv("ServerSave"),
		PlayerJoin:            os.Getenv("PlayerJoin"),
		PlayerLeave:           os.Getenv("PlayerLeave"),
	}
}
