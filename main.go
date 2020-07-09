package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"

	"github.com/maxsupermanhd/FactoCord-3.0/discord"
	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

var closing = false

func main() {
	fmt.Println("Welcome to FactoCord-3.0!")
	support.Config.MustLoad()

	discord.StartSession()

	go console()
	support.Factorio.Init(discord.ProcessFactorioLogLine)
	discord.Init()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, os.Kill)
	<-sc

	closing = true

	discord.Close()

	if support.Factorio.IsRunning() {
		fmt.Println("Waiting for factorio server to exit...")
		err := support.Factorio.Process.Wait()
		if support.Factorio.Process.ProcessState.Exited() {
			fmt.Println("\nFactorio server was closed, exit code", support.Factorio.Process.ProcessState.ExitCode())
		} else {
			fmt.Println("\nError waiting for factorio to exit")
			support.Panik(err, "Error waiting for factorio to exit")
		}
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
		support.Factorio.Send(string(line))
	}
}
