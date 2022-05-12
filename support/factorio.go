package support

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type FactorioLogWatcher struct {
	ProcessFunc func(string)
	buffer      string
}

func (t FactorioLogWatcher) Write(p []byte) (n int, err error) {
	t.buffer += string(p)
	lines := strings.Split(t.buffer, "\n")
	t.buffer = lines[len(lines)-1]
	for _, line := range lines[:len(lines)-1] {
		t.ProcessFunc(line)
	}
	return len(p), nil
}

func (t FactorioLogWatcher) Flush() {
	if t.buffer != "" {
		t.ProcessFunc(t.buffer)
		t.buffer = ""
	}
}

type factorioState struct {
	Process *exec.Cmd
	Pipe    *io.WriteCloser

	watcher       *io.Writer
	running       bool
	stopping      bool
	SaveRequested bool
	GameID        string
}

var Factorio factorioState

func (f *factorioState) Send(s string) bool {
	if f.Pipe == nil {
		return false
	}
	if s == "" {
		return false
	}
	if s[len(s)-1] != '\n' {
		s += "\n"
	}
	_, err := io.WriteString(*f.Pipe, s)
	Panik(err, "An error occurred when attempting send \""+s[:len(s)-1]+"\" to factorio")
	return err == nil
}

func (f *factorioState) IsRunning() bool {
	return f.running
}

func (f *factorioState) IsStopping() bool {
	return f.stopping
}

func (f *factorioState) Init(logger func(string)) {
	logging, err := os.OpenFile("factorio.log", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	Critical(err, "... when attempting to open factorio.log")

	factorioLogWatcher := FactorioLogWatcher{ProcessFunc: logger}
	tmpWatcher := io.MultiWriter(logging, os.Stdout, factorioLogWatcher)
	f.watcher = &tmpWatcher

	if Config.Autolaunch {
		factorioLogWatcher.Flush()
		Factorio.Start(nil)
	}
}

func (f *factorioState) Start(s *discordgo.Session) {
	if f.running {
		SendOptional(s, "The server is already running")
		return
	}
	if s != nil {
		SetTyping(s)
	}
	f.running = true
	f.Process = exec.Command(Config.Executable, Config.LaunchParameters...)
	f.Process.Stderr = os.Stderr
	f.Process.Stdout = *f.watcher
	pipe, err := f.Process.StdinPipe()
	Critical(err, "... when attempting to execute cmd.StdinPipe()")

	f.Pipe = &pipe

	err = f.Process.Start()
	Critical(err, "... when attempting to start the server")
}

func (f *factorioState) Stop(s *discordgo.Session) {
	if !f.running {
		SendOptional(s, "The server is already stopped")
		return
	}
	if f.stopping {
		SendOptional(s, "The server should be stopping")
		return
	}
	f.stopping = true
	f.Send("/quit")

	messageWaiting := SendOptional(s, "Waiting for factorio server to exit...")
	fmt.Println("Waiting for factorio server to exit...")
	err := f.Process.Wait()
	if f.Process.ProcessState.Exited() {
		if s != nil {
			messageWaiting.DeleteIfPassedLess(s, 10*time.Second)
			SendOptional(s, "Factorio server has **exited**")
		}
		fmt.Println("\nFactorio server was closed, exit code", f.Process.ProcessState.ExitCode())
	} else {
		fmt.Println("\nError waiting for factorio to exit")
		Panik(err, "Error waiting for factorio to exit")
	}
	f.Process = nil
	f.Pipe = nil
	f.running = false
	f.stopping = false
	f.SaveRequested = false
}

func FactorioVersion() (string, error) {
	cmd := exec.Command(Config.Executable, "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.Fields(string(out))[1], nil
}
