package main

import (
	"io"
	"log"
	"os"
)

func loggingSettings(logFile string) {
	logfile, _ := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	multiLogFile := io.MultiWriter(os.Stdout, logfile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(multiLogFile)
}

func main() {
	loggingSettings("log.txt")
	_, err := os.Open("undefined.txt")
	if err != nil {
		log.Fatalln("Exit", err)
	}

}
