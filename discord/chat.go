package discord

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/maxsupermanhd/FactoCord-3.0/commands"
	"github.com/maxsupermanhd/FactoCord-3.0/support"
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

	GuildChannel, err := Session.Channel(support.Config.FactorioChannelID)
	support.Critical(err, "... when attempting to read the Discord Guild")

	support.GuildID = GuildChannel.GuildID
}

func Init() {
	Session.AddHandler(messageCreate)
	Session.AddHandler(messageUpdate)
	// TODO add recover() ↑

	go CacheUpdater(Session)

	time.Sleep(3 * time.Second)
	err := Session.UpdateGameStatus(0, support.Config.GameName)
	support.Panik(err, "... when updating bot status")

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	support.SendMessage(Session, support.Config.Messages.BotStart)
}

func Close() {
	support.SendMessage(Session, support.Config.Messages.BotStop)

	// Cleanly close down the Discord session.
	err := Session.Close()
	support.Panik(err, "... when closing discord connection")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.ChannelID == support.Config.FactorioChannelID {
		support.MyLastMessage = false
		if strings.HasPrefix(m.Content, support.Config.Prefix) {
			input := strings.Replace(m.Content, support.Config.Prefix, "", 1)
			commands.RunCommand(input, s, m.Message)
			return
		}
		if regexp.MustCompile(fmt.Sprintf("^<@!?%s>", s.State.User.ID)).MatchString(m.Content) {
			// when bot is mentioned at the start of the message
			_, message := support.SplitAfter(m.Content, ">")
			message = strings.TrimSpace(message)
			if message == "" {
				support.SendFormat(s, "Current prefix is: `$`. You can use `$help` or ping me with a command")
			} else {
				commands.RunCommand(message, s, m.Message)
			}
			return
		}
		log.Print("[" + m.Author.Username + "] " + m.Content)
		// Pipes normal chat allowing it to be seen ingame
		if strings.TrimSpace(m.Content) != "" {
			// TODO? add color to mentions
			lines := strings.Split(m.ContentWithMentionsReplaced(), "\n")
			for i, line := range lines {
				if i != 0 {
					line = "[color=#6CFF3B]⬑[/color] " + line
				}
				lines[i] = fmt.Sprintf("<%s>: %s", colorUsername(m.Message), line)
				lines[i] = "[color=white]" + lines[i] + "[/color]"
				lines[i] = discordSignature + " " + lines[i]
			}
			support.Factorio.Send(strings.Join(lines, "\n"))
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
				attachmentType = "file: " + filename
			} else {
				attachmentType = fmt.Sprintf("image %dx%d", attachment.Width, attachment.Height)
			}
			attachmentType = fmt.Sprintf("[color=#35BFFF][%s][/color]", attachmentType)
			if strings.TrimSpace(m.Content) != "" {
				attachmentType = "[color=#6CFF3B]⬑[/color] " + attachmentType
			}
			message := fmt.Sprintf("[color=white]<%s>:[/color] %s", colorUsername(m.Message), attachmentType)
			support.Factorio.Send(discordSignature + " " + message)
		}
		return
	}
	if m.ChannelID == support.Config.FactorioConsoleChatID {
		fmt.Println("wrote to console from channel: \"", m.Content, "\"")
		support.SendTo(s, "wrote "+m.Content, support.Config.FactorioConsoleChatID)
		support.Factorio.Send(m.Content)
	}
	return
}

func messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	// TODO? refactor duplicate functions
	if m.Author == nil || m.Author.ID == s.State.User.ID || m.ChannelID != support.Config.FactorioChannelID {
		return
	}
	log.Print("[" + m.Author.Username + "]* " + m.Content)
	// Pipes normal chat allowing it to be seen ingame
	if strings.TrimSpace(m.Content) != "" {
		// TODO? add color to mentions
		lines := strings.Split(m.ContentWithMentionsReplaced(), "\n")
		for i, line := range lines {
			if i != 0 {
				line = "[color=#6CFF3B]⬑[/color] " + line
			}
			lines[i] = fmt.Sprintf("[color=#FFAA3B]<%s>*:[/color] %s", colorUsername(m.Message), line)
			lines[i] = "[color=white]" + lines[i] + "[/color]"
			lines[i] = discordSignature + " " + lines[i]
		}
		support.Factorio.Send(strings.Join(lines, "\n"))
	}
}

func colorUsername(message *discordgo.Message) string {
	if support.Config.IngameDiscordUserColors {
		color := Session.State.UserColor(message.Author.ID, message.ChannelID)
		if color == 0 { // some error
			return message.Author.Username
		} else {
			return fmt.Sprintf("[color=#%06x]%s[/color]", color, message.Author.Username)
		}
	} else {
		return message.Author.Username
	}
}
