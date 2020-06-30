package support

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"strings"
	"time"
)

func Send(s *discordgo.Session, message string) *discordgo.Message {
	sentMessage, err := s.ChannelMessageSend(Config.FactorioChannelID, message)
	if err != nil {
		Panik(err, "Failed to send message: "+message)
		return nil
	}
	return sentMessage
}

func SendOptional(s *discordgo.Session, message string) *discordgo.Message {
	if s == nil {
		return nil
	}
	return Send(s, message)
}

func SendMessage(s *discordgo.Session, message string) *discordgo.Message {
	if message != "" {
		return Send(s, message)
	}
	return nil
}

func SendFormat(s *discordgo.Session, message string) *discordgo.Message {
	return Send(s, FormatUsage(message))
}

func SendEmbed(s *discordgo.Session, embed *discordgo.MessageEmbed) *discordgo.Message {
	sentMessage, err := s.ChannelMessageSendEmbed(Config.FactorioChannelID, embed)
	if err != nil {
		Panik(err, fmt.Sprintf("Failed to send embed: %+v", embed))
		return nil
	}
	return sentMessage
}

func SendComplex(s *discordgo.Session, message *discordgo.MessageSend) *discordgo.Message {
	sentMessage, err := s.ChannelMessageSendComplex(Config.FactorioChannelID, message)
	if err != nil {
		Panik(err, fmt.Sprintf("Failed to send embed: %+v", message))
		return nil
	}
	return sentMessage
}

type MessageForDelete struct {
	ID      string
	Channel string
	Time    time.Time
}

func (m *MessageForDelete) Delete(s *discordgo.Session) {
	if m == nil || m.ID == "" {
		return
	}
	_ = s.ChannelMessageDelete(m.Channel, m.ID)
}
func (m *MessageForDelete) DeleteIfPassedLess(s *discordgo.Session, t time.Duration) {
	if m == nil || m.ID == "" {
		return
	}
	if time.Now().Before(m.Time.Add(t)) {
		m.Delete(s)
	}
}
func PrepareMessageDelete(m *discordgo.Message) *MessageForDelete {
	if m == nil {
		return &MessageForDelete{}
	}
	return &MessageForDelete{
		ID:      m.ID,
		Channel: m.ChannelID,
		Time:    time.Now(),
	}
}

func ChunkedMessageSend(s *discordgo.Session, message string) {
	lines := strings.Split(message, "\n")
	message = ""
	for _, line := range lines {
		if len(message)+len(line)+1 >= 2000 {
			Send(s, message)
			message = ""
		}
		message += "\n" + line
	}
	if len(message) > 0 {
		Send(s, message)
	}
}

func FormatUsage(s string) string {
	return strings.Replace(s, "$", Config.Prefix, -1)
}

func DeleteEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func SplitAt(s string, index int) (string, string) {
	if index < 0 {
		index += len(s)
	}
	return s[:index], s[index:]
}

func SplitBefore(s, sub string) (string, string) {
	index := strings.Index(s, sub)
	if index == -1 {
		return s, ""
	}
	return SplitAt(s, index)
}

func SplitAfter(s, sub string) (string, string) {
	index := strings.Index(s, sub)
	if index == -1 {
		return "", s
	}
	return SplitAt(s, index+len(sub))
}

func SplitDivide(s, sub string) (string, string) {
	index := strings.Index(s, sub)
	if index == -1 {
		return s, ""
	}
	return s[:index], s[index+len(sub):]
}

func QuoteSplit(s string, quote string) ([]string, bool) {
	var res []string
	firstQuote := -1
	for strings.Contains(s[firstQuote+len(quote):], quote) {
		if firstQuote == -1 {
			firstQuote = strings.Index(s, quote)
		} else {
			before := s[:firstQuote]
			if strings.TrimSpace(before) != "" {
				for _, x := range strings.Fields(before) {
					res = append(res, x)
				}
			}
			secondQuote := strings.Index(s[firstQuote+len(quote):], quote) + firstQuote + len(quote)
			unquoted := s[firstQuote+len(quote) : secondQuote]
			res = append(res, unquoted)
			s = s[secondQuote+len(quote):]
			firstQuote = -1
		}
	}
	mismatched := false
	if strings.TrimSpace(s) != "" {
		for _, x := range strings.Fields(s) {
			res = append(res, x)
			mismatched = mismatched || strings.Contains(x, quote)
		}
	}
	return res, mismatched
}

// FileExists checks if a file exists and is not a directory
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirExists checks if a directory exists and is not a file
func DirExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
