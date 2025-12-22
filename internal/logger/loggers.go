package logger

import (
	"fmt"
	"time"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

var globalReadline *readline.Instance

func SetReadline(rl *readline.Instance) {
	globalReadline = rl
}

func Info(message string) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	fmt.Println(color.YellowString(timestamp + " " + message))
	if globalReadline != nil {
		globalReadline.Refresh()
	}
}

func Error(message string) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	fmt.Println(color.RedString(timestamp + " " + message))
	if globalReadline != nil {
		globalReadline.Refresh()
	}
}

func Success(message string) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	fmt.Println(color.GreenString(timestamp + " " + message))
	if globalReadline != nil {
		globalReadline.Refresh()
	}
}

func Debug(message string) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	fmt.Println(color.CyanString(timestamp + " " + message))
	if globalReadline != nil {
		globalReadline.Refresh()
	}
}

func Warning(message string) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	fmt.Println(color.MagentaString(timestamp + " " + message))
	if globalReadline != nil {
		globalReadline.Refresh()
	}
}
