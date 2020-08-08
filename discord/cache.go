package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

func CacheUpdater(session *discordgo.Session) {
	// Wait 10 seconds on start up before continuing
	time.Sleep(10 * time.Second)

	for {
		count := CacheDiscordMembers(session)
		fmt.Printf("%s: discord members update: %d members\n", time.Now().Format("2006.01.02 15:04:05"), count)

		//sleep for 4 hours (caches every 4 hours)
		time.Sleep(4 * time.Hour)
	}
}

func CacheDiscordMembers(session *discordgo.Session) (count int) {
	after := ""
	limit := 1000

	for {
		members, err := session.GuildMembers(support.GuildID, after, limit)
		if err != nil {
			support.Panik(err, "... when requesting members")
			return
		}
		for _, member := range members {
			member.GuildID = support.GuildID
			err = session.State.MemberAdd(member)
			support.Panik(err, "... when adding member to state")
		}
		count += len(members)
		if len(members) < limit {
			break
		}
		after = members[len(members)-1].User.ID
	}
	return
}

// SearchForUser searches for the user to be mentioned.
func SearchForUser(name string) *discordgo.User {
	name = strings.Replace(name, "@", "", -1)
	guild, err := Session.State.Guild(support.GuildID)
	support.Panik(err, "... when getting guild")

	for _, member := range guild.Members {
		if strings.ToLower(member.Nick) == strings.ToLower(name) ||
			strings.ToLower(member.User.Username) == strings.ToLower(name) {
			return member.User
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
