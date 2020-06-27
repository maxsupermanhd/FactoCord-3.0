package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"./commands/admin"
	"./discord"
	"./support"
)

// Running is the boolean that tells if the server is running or not
var Running = false

var closing = false

// Factorio is a running factorio server instance
var Factorio *exec.Cmd

func main() {
	fmt.Println("Welcome to FactoCord-3.0!")
	support.Config.MustLoad()

	admin.R = &Running

	discord.StartSession()

	go factorioManager()
	go console()

	discord.Init()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, os.Kill)
	<-sc

	closing = true

	discord.Close()

	if Running {
		fmt.Println("Waiting for factorio server to exit...")
		err := Factorio.Wait()
		if Factorio.ProcessState.Exited() {
			fmt.Println("\nFactorio server was closed, exit code", Factorio.ProcessState.ExitCode())
		} else {
			fmt.Println("\nError waiting for factorio to exit")
			support.Panik(err, "Error waiting for factorio to exit")
		}
	}
}

func factorioManager() {
	logging, err := os.OpenFile("factorio.log", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	support.Critical(err, "... when attempting to open factorio.log")

	factorioLogWatcher := discord.FactorioLogWatcher{}
	mwriter := io.MultiWriter(logging, os.Stdout, factorioLogWatcher)
	for !closing {
		if !Running {
			Running = true
			factorioLogWatcher.Flush()
			Factorio = exec.Command(support.Config.Executable, support.Config.LaunchParameters...)
			Factorio.Stderr = os.Stderr
			Factorio.Stdout = mwriter
			var err error
			pipe, err := Factorio.StdinPipe()
			support.Critical(err, "... when attempting to execute cmd.StdinPipe()")

			support.FactorioPipe = &pipe

			err = Factorio.Start()
			support.Critical(err, "... when attempting to start the server")

			if admin.RestartCount > 0 {
				time.Sleep(3 * time.Second)
				support.Send(discord.Session, "Server restarted successfully!")
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func console() {
	Console := bufio.NewReader(os.Stdin)
	for !closing {
		line, _, err := Console.ReadLine()
		if err != nil {
			support.Panik(err, "An error occurred when attempting to read the input to pass as input to the console")
			return
		}
		support.SendToFactorio(string(line))
	}
}
