package discord

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"regexp"
	"strings"

	"../support"
)

type FactorioLogWatcher struct {
	buffer string
}

func (t FactorioLogWatcher) Write(p []byte) (n int, err error) {
	t.buffer += string(p)
	lines := strings.Split(t.buffer, "\n")
	t.buffer = lines[len(lines)-1]
	for _, line := range lines[:len(lines)-1] {
		ProcessFactorioLogLine(line)
	}
	return len(p), nil
}

func (t FactorioLogWatcher) Flush() {
	if t.buffer != "" {
		ProcessFactorioLogLine(t.buffer)
		t.buffer = ""
	}
}

var charRegexp = regexp.MustCompile("^\\d{4}[-/]\\d\\d[-/]\\d\\d \\d\\d:\\d\\d:\\d\\d ")
var factorioLogRegexp = regexp.MustCompile("^\\d+\\.\\d{3} ")

var forwardMessages = []*regexp.Regexp{
	regexp.MustCompile("^Player .+ doesn't exist."),
	regexp.MustCompile("^.+ wasn't banned."),
}

// ProcessFactorioLogLine pipes in-game chat to Discord.
func ProcessFactorioLogLine(line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}
	if charRegexp.FindString(line) != "" {
		line = line[len("0000-00-00 00:00:00 "):]
		processFactorioChat(strings.TrimSpace(line))
	} else if factorioLogRegexp.FindString(line) != "" {
		if strings.Contains(line, "Quitting: multiplayer error.") {
			support.SendMessage(Session, support.Config.Messages.ServerFail)
		}
		if strings.Contains(line, "Opening socket for broadcast") {
			support.SendMessage(Session, support.Config.Messages.ServerStart)
		}
		if strings.Contains(line, "Saving finished") {
			support.SendMessage(Session, "Saving finished!")
		}
		if strings.Contains(line, "Quitting multiplayer connection.") {
			support.SendMessage(Session, support.Config.Messages.ServerStop)
		}
	} else {
		for _, pattern := range forwardMessages {
			if pattern.FindString(line) != "" {
				support.Send(Session, line)
				return
			}
		}
	}
}

var chatStartRegexp = regexp.MustCompile("^\\[(CHAT|JOIN|LEAVE|KICK|BAN|DISCORD|DISCORD-EMBED)]")

func processFactorioChat(line string) {
	match := chatStartRegexp.FindStringSubmatch(line)
	if match == nil {
		return
	}
	messageType := match[1]
	integrationMessage := messageType == "DISCORD-EMBED" || messageType == "DISCORD"

	line = strings.TrimLeft(line[len(messageType)+2:], " ")
	if strings.HasPrefix(line, "<server>") {
		return
	}
	if messageType == "DISCORD" || messageType == "CHAT" {
		if strings.Contains(line, "@") {
			line = AddMentions(line)
		}
	}
	if support.Config.HaveServerEssentials {
		if messageType == "DISCORD-EMBED" {
			message := new(discordgo.MessageSend)
			err := json.Unmarshal([]byte(line), message)
			if err == nil {
				message.Tts = false
				support.SendComplex(Session, message)
			}
		} else if messageType == "DISCORD" {
			support.Send(Session, line)
		}
	} else if !integrationMessage {
		support.Send(Session, line)
	}
}
