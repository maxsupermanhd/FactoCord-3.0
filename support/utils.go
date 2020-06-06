package support

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

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

func ChunkedMessageSend(s *discordgo.Session, channel string, message string) {
	lines := strings.Split(message, "\n")
	message = ""
	for _, line := range lines {
		if len(message)+len(line)+1 >= 2000 {
			_, err := s.ChannelMessageSend(channel, message)
			if err != nil {
				fmt.Println("ChannelMessageSend failed")
				return
			}
			message = ""
		}
		message += "\n" + line
	}
	if len(message) > 0 {
		_, err := s.ChannelMessageSend(channel, message)
		if err != nil {
			fmt.Println("ChannelMessageSend failed")
			return
		}
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
