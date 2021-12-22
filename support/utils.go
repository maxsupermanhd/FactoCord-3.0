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
	chunks := ChunkMessage(message)
	for _, chunk := range chunks {
		Send(s, chunk)
	}
}

func Respond(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		Panik(err, "Failed to send message: "+content)
	}
}
func RespondChunked(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	chunks := ChunkMessage(content)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: chunks[0],
		},
	})
	if err != nil {
		Panik(err, "Failed to respond: "+content)
		return
	}
	chunks = chunks[1:]
	for _, chunk := range chunks {
		_, err := s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
			Content: chunk,
		})
		if err != nil {
			Panik(err, "Failed to send follow-up message: "+content)
			return
		}
	}
}
func RespondComplex(s *discordgo.Session, i *discordgo.InteractionCreate, data *discordgo.InteractionResponseData) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})
	if err != nil {
		Panik(err, fmt.Sprintf("Failed to send message: %v", data))
	}
}
func RespondDefer(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		Panik(err, "Failed to send message: "+content)
	}
}
func ResponseEdit(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	_, err := s.InteractionResponseEdit(s.State.User.ID, i.Interaction, &discordgo.WebhookEdit{
		Content: content,
	})
	if err != nil {
		Panik(err, "Failed to send message: "+content)
	}
}
func ResponseEditCompex(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
	_, err := s.InteractionResponseEdit(s.State.User.ID, i.Interaction, &discordgo.WebhookEdit{
		Embeds: []*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		Panik(err, fmt.Sprintf("Failed to send message: %v", embed))
	}
}

func ChunkMessage(message string) (res []string) {
	lines := strings.Split(message, "\n")
	message = ""
	for _, line := range lines {
		if len(message)+len(line)+1 >= 2000 {
			res = append(res, message)
			message = ""
		}
		message += "\n" + line
	}
	if len(message) > 0 {
		res = append(res, message)
	}
	return res
}

func SetTyping(s *discordgo.Session) {
	err := s.ChannelTyping(Config.FactorioChannelID)
	Panik(err, "... when sending 'typing' status")
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

func QuoteSpace(s string) string {
	if strings.ContainsRune(s, ' ') {
		s = "\"" + s + "\""
	}
	return s
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

func PluralS(x int) string {
	if x > 1 {
		return "s"
	}
	return ""
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
	Interaction               *discordgo.InteractionCreate
	Start, Progress, Finished string
}

func DownloadProgressUpdater(s *discordgo.Session, p *ProgressUpdate) {
	//message := p.Message
	//if message == nil {
	//	p.Message = Send(s, p.Start)
	//	message = p.Message
	//}
	time.Sleep(500 * time.Millisecond)
	for {
		if p.Error {
			return
		}
		if p.Transferred >= p.Total {
			break
		}
		percent := fmt.Sprintf("%2.1f", p.Percent())
		ResponseEdit(s, p.Interaction, FormatNamed(p.Progress, "percent", percent))
		time.Sleep(2 * time.Second)
	}
	ResponseEdit(s, p.Interaction, p.Finished)
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

func (l *TextListT) RenderWithoutHeading() string {
	if l.Error != "" {
		return l.Error
	}
	res := ""
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

type Command struct {
	Name        string
	Desc        string
	Usage       string
	Doc         string
	Admin       bool
	Command     func(s *discordgo.Session, i *discordgo.InteractionCreate)
	Subcommands []Command
	Options     []*discordgo.ApplicationCommandOption
	Choices     []*discordgo.ApplicationCommandOptionChoice
}

func (c *Command) Subcommand(s string) *Command {
	for _, subcommand := range c.Subcommands {
		if subcommand.Name == s {
			return &subcommand
		}
	}
	return nil
}

func (c *Command) ToSubcommand() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:         discordgo.ApplicationCommandOptionSubCommand,
		Name:         c.Name,
		Description:  c.Desc,
		Options:      c.Options,
		Autocomplete: false,
		Choices:      c.Choices,
	}
}

func (c *Command) ToCommand() *discordgo.ApplicationCommand {
	res := &discordgo.ApplicationCommand{
		Name:        c.Name,
		Description: c.Desc,
		Options:     c.Options,
	}
	for _, subcommand := range c.Subcommands {
		option := subcommand.ToSubcommand()
		res.Options = append(res.Options, option)
	}
	return res
}

func CompareOp(cmp int, op string) bool {
	switch op {
	case "=", "==":
		return cmp == 0
	case ">":
		return cmp == 1
	case ">=":
		return cmp >= 0
	case "<":
		return cmp == -1
	case "<=":
		return cmp <= 0
	}
	err := fmt.Errorf("`%s` is not a comparison operator", op)
	Panik(err, "")
	panic(err)
}
