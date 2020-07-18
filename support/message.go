package support

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

type MessageControlT struct {
	*discordgo.Message
	TimeSent time.Time
	Metadata string
}

var LastMessage *MessageControlT

var MyLastMessage bool

func MessageControl(m *discordgo.Message) *MessageControlT {
	if m == nil {
		return &MessageControlT{}
	}
	return &MessageControlT{
		Message:  m,
		TimeSent: time.Now(),
	}
}

func (m *MessageControlT) Edit(s *discordgo.Session, new string) *discordgo.Message {
	if m == nil || m.ID == "" {
		return nil
	}
	message, err := s.ChannelMessageEdit(m.ChannelID, m.ID, new)
	if err != nil {
		return nil
	}
	m.Message = message
	return message
}

func (m *MessageControlT) Delete(s *discordgo.Session) {
	if m == nil || m.ID == "" {
		return
	}
	_ = s.ChannelMessageDelete(m.ChannelID, m.ID)
	m.ID = ""
}

func (m *MessageControlT) DeleteIfPassedLess(s *discordgo.Session, t time.Duration) {
	if m == nil || m.ID == "" {
		return
	}
	if time.Now().Before(m.TimeSent.Add(t)) {
		m.Delete(s)
	}
}
