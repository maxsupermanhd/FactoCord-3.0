package discord

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"../commands"
	"../support"
)

var Session *discordgo.Session

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
		support.SendToFactorio(fmt.Sprintf("[Discord] <%s>: %s", m.Author.Username, strings.Replace(m.ContentWithMentionsReplaced(), "\n", fmt.Sprintf("\n[Discord] <%s>: ", m.Author.Username), -1)))
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

// Chat pipes in-game chat to Discord.
func ProcessFactorioLogLine(line string) {
	if strings.Contains(line, "Quitting: multiplayer error.") {
		support.Send(Session, support.Config.ServerFail)
	}
	if strings.Contains(line, "Info UDPSocket.cpp:39: Opening socket for broadcast") {
		support.Send(Session, support.Config.ServerStart)
	}
	if strings.Contains(line, "Info AppManagerStates.cpp:1843: Saving finished") {
		support.Send(Session, "Saving finished!")
	}
	if strings.Contains(line, "Info ServerMultiplayerManager.cpp:138: Quitting multiplayer connection.") {
		support.Send(Session, support.Config.ServerStop)
	}
	if support.Config.HaveServerEssentials == true {
		if strings.Contains(line, "[DISCORD]") ||
			strings.Contains(line, "[DISCORD-EMBED]") {
			if !strings.Contains(line, "<server>") || support.Config.PassConsoleChat {
				if strings.Contains(line, "[DISCORD-EMBED]") {
					TmpList := strings.Split(line, " ")
					message := new(discordgo.MessageSend)
					err := json.Unmarshal([]byte(strings.Join(TmpList[3:], " ")), message)
					if err == nil {
						message.Tts = false
						support.SendComplex(Session, message)
					}
				} else {
					TmpList := strings.Split(line, " ")
					TmpList[3] = strings.Replace(TmpList[3], ":", "", -1)
					if strings.Contains(strings.Join(TmpList, " "), "@") {
						index := support.LocateMentionPosition(TmpList)
						for _, position := range index {
							User := SearchForUser(TmpList[position])
							if User == nil {
								continue
							}
							TmpList[position] = User.Mention()
						}
					}
					support.Send(Session, strings.Join(TmpList[3:], " "))
				}
			}
		}
	} else {
		if strings.Contains(line, "[CHAT]") || strings.Contains(line, "[JOIN]") || strings.Contains(line, "[LEAVE]") || strings.Contains(line, "[KICK]") || strings.Contains(line, "[BAN]") {
			if !strings.Contains(line, "<server>") || support.Config.PassConsoleChat {
				if strings.Contains(line, "[JOIN]") ||
					strings.Contains(line, "[LEAVE]") {
					TmpList := strings.Split(line, " ")
					support.Send(Session, fmt.Sprintf("%s", strings.Join(TmpList[3:], " ")))
				} else {
					TmpList := strings.Split(line, " ")
					TmpList[3] = strings.Replace(TmpList[3], ":", "", -1)
					if strings.Contains(strings.Join(TmpList, " "), "@") {
						index := support.LocateMentionPosition(TmpList)
						for _, position := range index {
							User := SearchForUser(TmpList[position])
							if User == nil {
								continue
							}
							TmpList[position] = User.Mention()
						}
					}
					support.Send(Session, fmt.Sprintf("<%s>: %s", TmpList[3], strings.Join(TmpList[4:], " ")))
				}
			}
		}
	}

}
