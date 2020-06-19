package discord

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"../commands"
	"../support"
)

// fuck golang. it's shit
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

var Session *discordgo.Session

var discordSignature = "[color=#7289DA][Discord][/color]"

func StartSession() {
	fmt.Println("Starting bot..")

	var err error
	Session, err = discordgo.New("Bot " + support.Config.DiscordToken)
	support.Critical(err, "... when attempting to create the Discord session")

	err = Session.Open()
	support.Critical(err, "... when attempting to connect to Discord")
}

func Init() {
	Session.AddHandler(messageCreate)
	go CacheUpdater(Session)

	time.Sleep(3 * time.Second)
	err := Session.UpdateStatus(0, support.Config.GameName)
	support.Panik(err, "... when updating bot status")

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	if support.Config.SendBotStart {
		support.Send(Session, support.Config.BotStart)
	}
}

func Close() {
	if support.Config.BotStop != "" {
		support.Send(Session, support.Config.BotStop)
	}
	// Cleanly close down the Discord session.
	err := Session.Close()
	support.Panik(err, "... when closing discord connection")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.ChannelID == support.Config.FactorioChannelID {
		if strings.HasPrefix(m.Content, support.Config.Prefix) {
			input := strings.Replace(m.Content, support.Config.Prefix, "", 1)
			commands.RunCommand(input, s, m)
			return
		}
		log.Print("[" + m.Author.Username + "] " + m.Content)
		// Pipes normal chat allowing it to be seen ingame
		if strings.TrimSpace(m.Content) != "" {
			// TODO? add color to mentions
			lines := strings.Split(m.ContentWithMentionsReplaced(), "\n")
			for i, line := range lines {
				lines[i] = fmt.Sprintf("<%s>: %s", m.Author.Username, line)
				lines[i] = "[color=white]" + lines[i] + "[/color]"
				lines[i] = discordSignature + " " + lines[i]
			}
			support.SendToFactorio(strings.Join(lines, "\n"))
		}
		for _, attachment := range m.Attachments {
			attachmentType := ""
			if attachment.Width == 0 {
				filename := attachment.Filename
				if len(filename) > 20 {
					if strings.Contains(filename, ".") {
						dotIndex := strings.LastIndex(filename, ".")
						filename = filename[:min(dotIndex, 20)] + "..." + filename[dotIndex:]
					} else {
						filename = filename[:20] + "..."
					}
				}
				attachmentType = "файл: " + filename
			} else {
				attachmentType = fmt.Sprintf("изображение %dx%d", attachment.Width, attachment.Height)
			}
			message := fmt.Sprintf("[color=white]<%s>:[/color] [color=#35BFFF][%s][/color]", m.Author.Username, attachmentType)
			support.SendToFactorio(discordSignature + " " + message)
		}
		return
	}
	if m.ChannelID == support.Config.FactorioConsoleChatID {
		fmt.Println("wrote to console from channel: \"", m.Content, "\"")
		support.Send(s, "wrote "+m.Content)
		support.SendToFactorio(m.Content)
	}
	return
}

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
		if support.Config.PassConsoleChat {
			line = line[len("0000-00-00 00:00:00 "):]
			processFactorioChat(strings.TrimSpace(line))
		}
	} else if factorioLogRegexp.FindString(line) != "" {
		if strings.Contains(line, "Quitting: multiplayer error.") {
			support.Send(Session, support.Config.ServerFail)
		}
		if strings.Contains(line, "Opening socket for broadcast") {
			support.Send(Session, support.Config.ServerStart)
		}
		if strings.Contains(line, "Saving finished") {
			support.Send(Session, "Saving finished!")
		}
		if strings.Contains(line, "Quitting multiplayer connection.") {
			support.Send(Session, support.Config.ServerStop)
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
