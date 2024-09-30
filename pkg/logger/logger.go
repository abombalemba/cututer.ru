package logger

import (
	"os"
	"log"

	"cututer/tools"
)

var (
	logger *log.Logger
	fileLog *os.File
)

func CreateLogger() {
	filename := tools.GetNow()

	fileLog, err := os.OpenFile("../../logs/" + filename + ".log", os.O_CREATE | os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("Error opening log file: ", err)
		return
	}

	logger = log.New(fileLog, "", log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}

func GetLogger() *log.Logger {
	return logger
}

func CloseLogger() {
	if err := fileLog.Close(); err != nil {
		log.Fatalln("Error closing log file: ", err)
	}
}