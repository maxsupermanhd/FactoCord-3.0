package discord

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"

	"../support"
)

func CacheUpdater(session *discordgo.Session) {
	// Wait 10 seconds on start up before continuing
	time.Sleep(10 * time.Second)

	for {
		CacheDiscordMembers(session)
		//sleep for 4 hours (caches every 4 hours)
		time.Sleep(4 * time.Hour)
	}
}

// UserList is a struct for member info.
type UserList struct {
	UserID string
	Nick   string
	User   *discordgo.User
}

// Users is a slice of UserList.
var Users []UserList

// CacheDiscordMembers caches the users list to be searched.
func CacheDiscordMembers(s *discordgo.Session) {
	// Clear the users list
	Users = nil

	GuildChannel, err := s.Channel(support.Config.FactorioChannelID)
	support.Panik(err, "... when attempting to read the Discord Guild")

	GuildID := GuildChannel.GuildID
	members, err := s.State.Guild(GuildID)
	support.Panik(err, "... when attempting to read the Discord Guild Members")

	for _, member := range members.Members {
		Users = append(Users, UserList{UserID: member.User.ID, Nick: member.Nick,
			User: member.User})
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

func AddMentions(message string) string {
	if !strings.Contains(message, "@") {
		return message
	}
	words := strings.Split(message, " ")
	for i, word := range words {
		if len(word) >= 2 && word[0] == '@' {
			User := SearchForUser(word[1:])
			if User == nil {
				continue
			}
			words[i] = User.Mention()
		}
	}
	message = strings.Join(words, " ")
	return message
}
