package support

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Send(s *discordgo.Session, message string) {
	_, err := s.ChannelMessageSend(Config.FactorioChannelID, message)
	if err != nil {
		Panik(err, "Failed to send message: "+message)
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

// SearchForUser searches for the user to be mentioned.
func SearchForUser(name string) *discordgo.User {
	name = strings.Replace(name, "@", "", -1)
	for _, user := range Users {
		if strings.ToLower(user.Nick) == strings.ToLower(name) ||
			strings.ToLower(user.User.Username) == strings.ToLower(name) {
			return user.User
		}
	}
	return nil
}

// LocateMentionPosition locates the position in a string list for the discord mention.
func LocateMentionPosition(List []string) []int {
	positionlist := []int{}
	for i, String := range List {
		if strings.Contains(String, "@") {
			positionlist = append(positionlist, i)
		}
	}
	return positionlist
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
