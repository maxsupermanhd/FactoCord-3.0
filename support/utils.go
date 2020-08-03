package support

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"strings"
	"time"
)

func Send(s *discordgo.Session, message string) *MessageControlT {
	sentMessage, err := s.ChannelMessageSend(Config.FactorioChannelID, message)
	if err != nil {
		Panik(err, "Failed to send message: "+message)
		return nil
	}
	LastMessage = MessageControl(sentMessage)
	MyLastMessage = true
	return LastMessage
}

func SendOptional(s *discordgo.Session, message string) *MessageControlT {
	if s == nil {
		return nil
	}
	return Send(s, message)
}

func SendMessage(s *discordgo.Session, message string) *MessageControlT {
	if message != "" {
		return Send(s, message)
	}
	return nil
}

func SendFormat(s *discordgo.Session, message string) *MessageControlT {
	return Send(s, FormatUsage(message))
}

func SendEmbed(s *discordgo.Session, embed *discordgo.MessageEmbed) *MessageControlT {
	sentMessage, err := s.ChannelMessageSendEmbed(Config.FactorioChannelID, embed)
	if err != nil {
		Panik(err, fmt.Sprintf("Failed to send embed: %+v", embed))
		return nil
	}
	LastMessage = MessageControl(sentMessage)
	MyLastMessage = true
	return LastMessage
}

func SendComplex(s *discordgo.Session, message *discordgo.MessageSend) *MessageControlT {
	sentMessage, err := s.ChannelMessageSendComplex(Config.FactorioChannelID, message)
	if err != nil {
		Panik(err, fmt.Sprintf("Failed to send embed: %+v", message))
		return nil
	}
	LastMessage = MessageControl(sentMessage)
	MyLastMessage = true
	return LastMessage
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

func FormatNamed(format, name, value string) string {
	return strings.Replace(format, "{"+name+"}", value, 1)
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

func Unique(strs []string) []string {
	s := make([]string, len(strs))
	copy(s, strs)
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i] == s[j] {
				copy(s[j:], s[j+1:])
				s = s[:len(s)-1]
			}
		}
	}
	return s
}

func UniqueFunc(objs []interface{}, f func(interface{}, interface{}) bool) []interface{} {
	o := make([]interface{}, len(objs))
	copy(o, objs)
	for i := 0; i < len(o); i++ {
		for j := i + 1; j < len(o); j++ {
			if f(o[i], o[j]) {
				copy(o[j:], o[j+1:])
				o = o[:len(o)-1]
			}
		}
	}
	return o
}

func IsUnique(s []string) bool {
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i] == s[j] {
				return false
			}
		}
	}
	return true
}

func AnyTwo(o []interface{}, f func(interface{}, interface{}) bool) bool {
	for i := 0; i < len(o); i++ {
		for j := i + 1; j < len(o); j++ {
			if f(o[i], o[j]) {
				return true
			}
		}
	}
	return false
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

type WriteCounter struct {
	Total       uint64
	Transferred uint64
	Error       bool
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Transferred += uint64(n)
	return n, nil
}

func (wc *WriteCounter) Percent() float32 {
	return float32(wc.Transferred) * 100 / float32(wc.Total)
}

type ProgressUpdate struct {
	*WriteCounter
	Message                   *MessageControlT
	Start, Progress, Finished string
}

func DownloadProgressUpdater(s *discordgo.Session, p *ProgressUpdate) {
	message := p.Message
	if message == nil {
		p.Message = Send(s, p.Start)
		message = p.Message
	}
	time.Sleep(500 * time.Millisecond)
	for {
		if p.Error {
			return
		}
		if p.Transferred >= p.Total {
			break
		}
		percent := fmt.Sprintf("%2.1f", p.Percent())
		message.Edit(s, FormatNamed(p.Progress, "percent", percent))
		time.Sleep(2 * time.Second)
	}
	message.Edit(s, p.Finished)
}

type TextListT struct {
	Heading     string
	List        []string
	Indentation string
	None        string
	Error       string
}

func DefaultTextList(heading string) TextListT {
	return TextListT{
		Heading:     heading,
		List:        []string{},
		Indentation: "    ",
		None:        " **None**",
	}
}

func (l *TextListT) IsEmpty() bool {
	return len(l.List) == 0 && l.Error == ""
}

func (l *TextListT) NotEmpty() bool {
	return !l.IsEmpty()
}

func (l *TextListT) Len() int {
	return len(l.List)
}

func (l *TextListT) Append(s string) {
	l.List = append(l.List, s)
}

func (l *TextListT) AddToLast(s string) {
	l.List[len(l.List)-1] += s
}

func (l *TextListT) FormatHeaderWithLength() {
	l.Heading = fmt.Sprintf(l.Heading, len(l.List))
}

func (l *TextListT) Render() string {
	if l.Error != "" {
		return l.Error
	}
	res := l.Heading
	if len(l.List) == 0 {
		res += l.None
	} else {
		for _, x := range l.List {
			res += "\n" + l.Indentation + x
		}
	}
	return res
}

func (l *TextListT) RenderNotEmpty() string {
	if l.IsEmpty() {
		return ""
	} else {
		return l.Render()
	}
}
