package support

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

// Exit exits the application.
func Exit(ExitCode int) {
	exit, rip := os.OpenFile(".exit", os.O_RDWR|os.O_CREATE, 0666)
	if rip != nil {
		fmt.Println("Error opening \".exit\" file")
		panic(rip)
	}
	_, rip = exit.WriteString(fmt.Sprintf("%d", ExitCode))
	if rip != nil {
		fmt.Println("Error writing to \".exit\"")
		panic(rip)
	}
	os.Exit(ExitCode)
}

// Panik checks error != nil and logs the error without exiting the app
// P.S. it's a meme
func Panik(err error, message string) {
	if err == nil {
		return
	}
	if strings.HasPrefix(message, "...") {
		message = "An error occurred " + strings.TrimSpace(message[3:])
	}

	res := fmt.Sprintf("%s", time.Now())
	if pc, fn, line, ok := runtime.Caller(1); ok {
		res += fmt.Sprintf(", %s @ %s:%d", runtime.FuncForPC(pc).Name(), fn, line)
	}
	if message == "" {
		res += fmt.Sprintf("\n\t%v\n", err)
	} else {
		res += fmt.Sprintf("\n\t%s: %v\n", message, err)
	}

	errorLog, rip := os.OpenFile("error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// If we encounter an error here, something is seriously wrong.
	if rip != nil {
		fmt.Println("Error occurred:")
		fmt.Println(res)
		fmt.Println("But error.log was inaccessible")
		panic(rip)
	}
	defer errorLog.Close()
	_, err = errorLog.WriteString(res)
	if err != nil {
		fmt.Println("Error occurred:")
		fmt.Println(res)
		fmt.Println("But error.log could not be written to")
		panic(rip)
	}

	if message == "" {
		fmt.Println("Oops, it looks like an error happened!")
	} else {
		fmt.Println(message)
	}
	fmt.Println("You can post your issue on https://github.com/maxsupermanhd/FactoCord-3.0/issues")
}

// Critical checks error != nil, logs the error and closes the app
func Critical(err error, message string) {
	if err != nil {
		Panik(err, message)
		Exit(1)
	}
}
