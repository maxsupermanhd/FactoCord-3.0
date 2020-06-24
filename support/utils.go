package support

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Send(s *discordgo.Session, message string) {
	_, err := s.ChannelMessageSend(Config.FactorioChannelID, message)
	if err != nil {
		Panik(err, "Failed to send message: "+message)
	}
}

func SendMessage(s *discordgo.Session, message string) {
	if message != "" {
		Send(s, message)
	}
}

func SendFormat(s *discordgo.Session, message string) {
	Send(s, FormatUsage(message))
}

func SendEmbed(s *discordgo.Session, embed *discordgo.MessageEmbed) {
	_, err := s.ChannelMessageSendEmbed(Config.FactorioChannelID, embed)
	if err != nil {
		Panik(err, fmt.Sprintf("Failed to send embed: %+v", embed))
	}
}

func SendComplex(s *discordgo.Session, message *discordgo.MessageSend) {
	_, err := s.ChannelMessageSendComplex(Config.FactorioChannelID, message)
	if err != nil {
		Panik(err, fmt.Sprintf("Failed to send embed: %+v", message))
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

// FileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
