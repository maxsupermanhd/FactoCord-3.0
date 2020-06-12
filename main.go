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
	"time"

	"./commands"
	"./commands/admin"
	"./support"
	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

// Running is the boolean that tells if the server is running or not
var Running bool

var Close bool = false

// Factorio is a running factorio server instance
var Factorio *exec.Cmd

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

	logging, err := os.OpenFile("factorio.log", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	support.Critical(err, "... when attempting to open factorio.log")

	mwriter := io.MultiWriter(logging, os.Stdout)
	var wg sync.WaitGroup
	wg.Add(3)

	go factorioManager(&wg, &mwriter)
	go console()
	go discordCache()

	discord()
}

func factorioManager(wg *sync.WaitGroup, mwriter *io.Writer) {
	defer wg.Done()
	for !Close {
		// If the process is already running DO NOT RUN IT AGAIN
		if !Running {
			Running = true
			Factorio = exec.Command(support.Config.Executable, support.Config.LaunchParameters...)
			Factorio.Stderr = os.Stderr
			Factorio.Stdout = *mwriter
			var err error
			Pipe, err = Factorio.StdinPipe()
			support.Critical(err, "... when attempting to execute cmd.StdinPipe()")

			err = Factorio.Start()
			support.Critical(err, "... when attempting to start the server")

			if admin.RestartCount > 0 {
				time.Sleep(3 * time.Second)
				support.Send(Session, "Server restarted successfully!")
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func discord() {
	// No hard coding the token }:<
	discordToken := support.Config.DiscordToken

	admin.P = &Pipe
	fmt.Println("Starting bot..")
	bot, err := discordgo.New("Bot " + discordToken)
	support.Critical(err, "... when attempting to create the Discord session")
	Session = bot

	err = bot.Open()
	support.Critical(err, "... when attempting to connect to Discord")

	bot.AddHandler(messageCreate)
	go support.Chat(bot)

	time.Sleep(3 * time.Second)
	err = bot.UpdateStatus(0, support.Config.GameName)
	support.Panik(err, "... when updating bot status")

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	if support.Config.SendBotStart {
		support.Send(bot, support.Config.BotStart)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, os.Kill)
	<-sc

	Close = true

	if support.Config.BotStop != "" {
		support.Send(bot, support.Config.BotStop)
	}
	// Cleanly close down the Discord session.
	err = bot.Close()
	support.Panik(err, "... when closing discord connection")

	if Running {
		fmt.Println("Waiting for factorio server to exit...")
		err = Factorio.Wait()
		if Factorio.ProcessState.Exited() {
			fmt.Println("\nFactorio server was closed, exit code", Factorio.ProcessState.ExitCode())
		} else {
			fmt.Println("\nError waiting for factorio to exit")
			support.Panik(err, "Error waiting for factorio to exit")
		}
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.ChannelID == support.Config.FactorioChannelID {
		if strings.HasPrefix(m.Content, support.Config.Prefix) {
			input := strings.Replace(m.Content, support.Config.Prefix, "", 1)
			commands.RunCommand(input, s, m)
			return
		}
		log.Print("[" + m.Author.Username + "] " + m.Content)
		// Pipes normal chat allowing it to be seen ingame
		_, err := io.WriteString(Pipe, fmt.Sprintf("[Discord] <%s>: %s\r\n", m.Author.Username, strings.Replace(m.ContentWithMentionsReplaced(), "\n", fmt.Sprintf("\n[Discord] <%s>: ", m.Author.Username), -1)))
		support.Panik(err, "An error occurred when attempting to pass Discord chat to in-game")
		return
	}
	if m.ChannelID == support.Config.FactorioConsoleChatID {
		fmt.Println("wrote to console from channel: \"", fmt.Sprintf("%s", m.Content), "\"")
		support.Send(s, fmt.Sprintf("wrote %s", m.Content))
		_, err := io.WriteString(Pipe, fmt.Sprintf("%s\n", m.Content))
		support.Panik(err, "An error occurred when attempting to pass Discord console to in-game")
	}
	return
}

func discordCache() {
	// Wait 10 seconds on start up before continuing
	time.Sleep(10 * time.Second)

	for !Close {
		support.CacheDiscordMembers(Session)
		//sleep for 4 hours (caches every 4 hours)
		time.Sleep(4 * time.Hour)
	}
}

func console() {
	Console := bufio.NewReader(os.Stdin)
	for !Close {
		line, _, err := Console.ReadLine()
		if err != nil {
			support.Panik(err, "An error occurred when attempting to read the input to pass as input to the console")
			return
		}
		_, err = io.WriteString(Pipe, fmt.Sprintf("%s\n", line))
		if err != nil {
			support.Panik(err, "An error occurred when attempting to pass input to the console")
			return
		}
	}
}
