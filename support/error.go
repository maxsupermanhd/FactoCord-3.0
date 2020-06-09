package support

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

// ErrorLog logs the error to file and then exits the application with an
// exit code of 1.
func ErrorLog(err error) {
	errorlog, rip := os.OpenFile("error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// If we encounter an error here, something is seriously wrong.
	if rip != nil {
		panic(rip)
	}
	defer errorlog.Close()
	errorlog.WriteString(fmt.Sprintf("%s\n", err))
	fmt.Println("Opps, it looks like an error happened!")
	fmt.Println("Please post your error.log on https://github.com/maxsupermanhd/FactoCord-3.0/issues")
	Exit(1)
}

func Panik(err error, message string) {
	if err == nil {
		return
	}
	if message == "" {
		message = "An error occurred"
	} else if strings.HasPrefix(message, "...") {
		message = "An error occurred " + strings.TrimSpace(message[3:])
	}

	res := fmt.Sprintf("%s", time.Now())
	if pc, fn, line, ok := runtime.Caller(1); ok {
		res += fmt.Sprintf(", %s @ %s:%d", runtime.FuncForPC(pc).Name(), fn, line)
	}
	res += fmt.Sprintf("\n\t%s: %v\n", message, err)

	errorLog, rip := os.OpenFile("error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// If we encounter an error here, something is seriously wrong.
	if rip != nil {
		fmt.Println("Error occured:\n")
		fmt.Println(res)
		fmt.Println("But error.log was unaccessible")
		panic(rip)
	}
	defer errorLog.Close()
	_, err = errorLog.WriteString(res)
	if err != nil {
		fmt.Println("Error occured:\n")
		fmt.Println(res)
		fmt.Println("But error.log could not be written to")
		panic(rip)
	}

	fmt.Println("Opps, it looks like an error happened!")
	fmt.Println("Please post your issue on https://github.com/maxsupermanhd/FactoCord-3.0/issues")
}

func Critical(err error, message string) {
	if err != nil {
		Panik(err, message)
		Exit(1)
	}
}

// P.S. panik is a meme
