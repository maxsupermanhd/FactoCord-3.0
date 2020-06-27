package support

import "io"

var FactorioPipe *io.WriteCloser

func SendToFactorio(s string) bool {
	if FactorioPipe == nil {
		return false
	}
	if s[len(s)-1] != '\n' {
		s += "\n"
	}
	_, err := io.WriteString(*FactorioPipe, s)
	Panik(err, "An error occurred when attempting send \""+s[:len(s)-1]+"\" to factorio")
	return err == nil
}
