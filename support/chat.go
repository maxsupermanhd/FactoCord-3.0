package support

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/hpcloud/tail"
)

// Chat pipes in-game chat to Discord.
func Chat(s *discordgo.Session, m *discordgo.MessageCreate) {
	for {
		t, err := tail.TailFile("factorio.log", tail.Config{Follow: true})
		if err != nil {
			ErrorLog(fmt.Errorf("%s: An error occurred when attempting to tail factorio.log\nDetails: %s", time.Now(), err))
		}
		for line := range t.Lines {
			if strings.Contains(line.Text, "Quitting: multiplayer error.") {
				s.ChannelMessageSend(Config.FactorioChannelID, Config.ServerFail)
			}
			if strings.Contains(line.Text, "Info UDPSocket.cpp:39: Opening socket for broadcast") {
				s.ChannelMessageSend(Config.FactorioChannelID, Config.ServerStart)
			}
			if strings.Contains(line.Text, "Info AppManagerStates.cpp:1843: Saving finished") {
				s.ChannelMessageSend(Config.FactorioChannelID, "Saving finished!")
			}
			if strings.Contains(line.Text, "Info ServerMultiplayerManager.cpp:138: Quitting multiplayer connection.") {
				s.ChannelMessageSend(Config.FactorioChannelID, Config.ServerStop)
			}
			if Config.HaveServerEssentials == true {
				if strings.Contains(line.Text, "[DISCORD]") ||
					strings.Contains(line.Text, "[DISCORD-EMBED]") {
					if !strings.Contains(line.Text, "<server>") || Config.PassConsoleChat {
						if strings.Contains(line.Text, "[DISCORD-EMBED]") {
							TmpList := strings.Split(line.Text, " ")
							message := new(discordgo.MessageSend)
							err := json.Unmarshal([]byte(fmt.Sprintf("%s", strings.Join(TmpList[3:], " "))), message)
							if err == nil {
								message.Tts = false
								_, serrn := s.ChannelMessageSendComplex(Config.FactorioChannelID, message)
								if serrn != nil {
									fmt.Println(fmt.Sprintf("Discord-embed error %s\n", serrn))
								}
							}
						} else {
							TmpList := strings.Split(line.Text, " ")
							TmpList[3] = strings.Replace(TmpList[3], ":", "", -1)
							if strings.Contains(strings.Join(TmpList, " "), "@") {
								index := LocateMentionPosition(TmpList)
								for _, position := range index {
									User := SearchForUser(TmpList[position])
									if User == nil {
										continue
									}
									TmpList[position] = User.Mention()
								}
							}
							s.ChannelMessageSend(Config.FactorioChannelID, fmt.Sprintf("%s", strings.Join(TmpList[3:], " ")))
						}
					}
				}
			} else {
				if strings.Contains(line.Text, "[CHAT]") || strings.Contains(line.Text, "[JOIN]") || strings.Contains(line.Text, "[LEAVE]") || strings.Contains(line.Text, "[KICK]") || strings.Contains(line.Text, "[BAN]") {
					if !strings.Contains(line.Text, "<server>") || Config.PassConsoleChat {
						if strings.Contains(line.Text, "[JOIN]") ||
							strings.Contains(line.Text, "[LEAVE]") {
							TmpList := strings.Split(line.Text, " ")
							s.ChannelMessageSend(Config.FactorioChannelID, fmt.Sprintf("%s", strings.Join(TmpList[3:], " ")))
						} else {
							TmpList := strings.Split(line.Text, " ")
							TmpList[3] = strings.Replace(TmpList[3], ":", "", -1)
							if strings.Contains(strings.Join(TmpList, " "), "@") {
								index := LocateMentionPosition(TmpList)
								for _, position := range index {
									User := SearchForUser(TmpList[position])
									if User == nil {
										continue
									}
									TmpList[position] = User.Mention()
								}
							}
							s.ChannelMessageSend(Config.FactorioChannelID, fmt.Sprintf("<%s>: %s", TmpList[3], strings.Join(TmpList[4:], " ")))
						}
					}
				}
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
}
