package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"./support"
	"./commands"
	"./commands/admin"
	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

// Running is the boolean that tells if the server is running or not
var Running bool

var SaveRequested bool = false

// Pipe is an WriteCloser interface
var Pipe io.WriteCloser

// Session is a discordgo session
var Session *discordgo.Session

func main() {
	fmt.Println("Welcome to FactoCord-3.0!")
	support.Config.LoadEnv()
	Running = false
	admin.R = &Running

	// Do not exit the app on this error.
	if err := os.Remove("factorio.log"); err != nil {
		fmt.Println("Factorio.log doesn't exist, continuing anyway")
	}

	logging, err := os.OpenFile("factorio.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		support.ErrorLog(fmt.Errorf("%s: An error occurred when attempting to open factorio.log\nDetails: %s", time.Now(), err))
	}

	mwriter := io.MultiWriter(logging, os.Stdout)
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		for {
			// If the process is already running DO NOT RUN IT AGAIN
			if !Running {
				Running = true
				cmd := exec.Command(support.Config.Executable, support.Config.LaunchParameters...)
				cmd.Stderr = os.Stderr
				cmd.Stdout = mwriter
				Pipe, err = cmd.StdinPipe()
				if err != nil {
					support.ErrorLog(fmt.Errorf("%s: An error occurred when attempting to execute cmd.StdinPipe()\nDetails: %s", time.Now(), err))
				}

				err := cmd.Start()

				if err != nil {
					support.ErrorLog(fmt.Errorf("%s: An error occurred when attempting to start the server\nDetails: %s", time.Now(), err))
				}
				if admin.RestartCount > 0 {
					time.Sleep(3 * time.Second)
					Session.ChannelMessageSend(support.Config.FactorioChannelID,
						"Server restarted successfully!")
				}
			}
			time.Sleep(5 * time.Second)
		}
	}()

	go func() {
		Console := bufio.NewReader(os.Stdin)
		for {
			line, _, err := Console.ReadLine()
			if err != nil {
				support.ErrorLog(fmt.Errorf("%s: An error occurred when attempting to read the input to pass as input to the console\nDetails: %s", time.Now(), err))
			}
			_, err = io.WriteString(Pipe, fmt.Sprintf("%s\n", line))
			if err != nil {
				support.ErrorLog(fmt.Errorf("%s: An error occurred when attempting to pass input to the console\nDetails: %s", time.Now(), err))
			}
		}
	}()

	go func() {
		// Wait 10 seconds on start up before continuing
		time.Sleep(10 * time.Second)

		for {
			support.CacheDiscordMembers(Session)
			//sleep for 4 hours (caches every 4 hours)
			time.Sleep(4 * time.Hour)
		}
	}()
	discord()
}

func discord() {
	// No hard coding the token }:<
	discordToken := support.Config.DiscordToken
	//commands.RegisterCommands()
	admin.P = &Pipe
	fmt.Println("Starting bot..")
	bot, err := discordgo.New("Bot " + discordToken)
	Session = bot
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		support.ErrorLog(fmt.Errorf("%s: An error occurred when attempting to create the Discord session\nDetails: %s", time.Now(), err))
		return
	}

	err = bot.Open()

	if err != nil {
		fmt.Println("error opening connection,", err)
		support.ErrorLog(fmt.Errorf("%s: An error occurred when attempting to connect to Discord\nDetails: %s", time.Now(), err))
		return
	}

	bot.AddHandler(messageCreate)
	bot.AddHandlerOnce(support.Chat)
	time.Sleep(3 * time.Second)
	bot.UpdateStatus(0, support.Config.GameName)
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	if support.Config.SendBotStart {
		bot.ChannelMessageSend(support.Config.FactorioChannelID, support.Config.BotStart)
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	bot.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.ChannelID == support.Config.FactorioChannelID {
		if strings.HasPrefix(m.Content, support.Config.Prefix) {
			//command := strings.Split(m.Content[1:len(m.Content)], " ")
			//name := strings.ToLower(command[0])
			input := strings.Replace(m.Content, support.Config.Prefix, "", -1)
			commands.RunCommand(input, s, m)
			return
		}
		log.Print("[" + m.Author.Username + "] " + m.Content)
		// Pipes normal chat allowing it to be seen ingame
		_, err := io.WriteString(Pipe, fmt.Sprintf("[Discord] <%s>: %s\r\n", m.Author.Username, strings.Replace(m.ContentWithMentionsReplaced(), "\n", fmt.Sprintf("\n[Discord] <%s>: ", m.Author.Username), -1)))
		if err != nil {
			support.ErrorLog(fmt.Errorf("%s: An error occurred when attempting to pass Discord chat to in-game\nDetails: %s", time.Now(), err))
		}
		return
	}
	if m.ChannelID == support.Config.FactorioConsoleChatID {
		fmt.Println("wrote to console from channel: \"", fmt.Sprintf("%s", m.Content), "\"")
		s.ChannelMessageSend(support.Config.FactorioConsoleChatID, fmt.Sprintf("wrote %s", m.Content))
		_, err := io.WriteString(Pipe, fmt.Sprintf("%s\n", m.Content))
		if err != nil {
			support.ErrorLog(fmt.Errorf("%s: An error occurred when attempting to pass Discord console to in-game\nDetails: %s", time.Now(), err))
		}
	}
	return
}

func CheckAdmin(ID string) bool {
	for _, admin := range support.Config.AdminIDs {
		if ID == admin {
			return true
		}
	}
	return false
}
