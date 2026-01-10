package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/chzyer/readline"
)

// bonus: colored output without external dependencies

var globalReadline *readline.Instance
var logFile *os.File

func InitializeLogFile() {
	var err error
	logFile, err = os.Create("taskMaster.log")
	if err != nil {
		fmt.Println("Failed to create log file:", err)
		panic(err)
	}
}

func SetReadline(rl *readline.Instance) {
	globalReadline = rl
}

func Info(message string) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	_, err := logFile.WriteString(timestamp + " " + message + "\n")
	if err != nil {
		fmt.Println("Failed to write to log file:", err)
	}
	// fmt.Println("\033[33m" + timestamp + " " + message + "\033[0m")
	if globalReadline != nil {
		globalReadline.Refresh()
	}
}

func Error(message string) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	_, err := logFile.WriteString(timestamp + " " + message + "\n")
	if err != nil {
		fmt.Println("Failed to write to log file:", err)
	}
	fmt.Println("\033[31m" + timestamp + " " + message + "\033[0m")
	if globalReadline != nil {
		globalReadline.Refresh()
	}
}

func Success(message string) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	_, err := logFile.WriteString(timestamp + " " + message + "\n")
	if err != nil {
		fmt.Println("Failed to write to log file:", err)
	}
	// fmt.Println(color.GreenString(timestamp + " " + message))
	if globalReadline != nil {
		globalReadline.Refresh()
	}
}

func Debug(message string) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	_, err := logFile.WriteString(timestamp + " " + message + "\n")
	if err != nil {
		fmt.Println("Failed to write to log file:", err)
	}
	fmt.Println("\033[36m" + timestamp + " " + message + "\033[0m")
	if globalReadline != nil {
		globalReadline.Refresh()
	}
}

func Warning(message string) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	_, err := logFile.WriteString(timestamp + " " + message + "\n")
	if err != nil {
		fmt.Println("Failed to write to log file:", err)
	}
	fmt.Println("\033[35m" + timestamp + " " + message + "\033[0m")
	if globalReadline != nil {
		globalReadline.Refresh()
	}
}
