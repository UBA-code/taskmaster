package main

import (
	"io"
	"os"
	"strings"

	"github.com/uba-code/taskmaster/internal/cli"
	"github.com/uba-code/taskmaster/internal/config"
	"github.com/uba-code/taskmaster/internal/logger"
)

func main() {
	if len(os.Args) < 2 {
		panic("Usage: taskmaster <config-file>")
	}
	var cfg = config.ParseConfig(os.Args[1])
	_ = cfg

	var rl = cli.ReadlineInit()
	defer rl.Close()
	logger.SetReadline(rl)
	var tasks = cli.NewTasksObj(cfg)

	//* cli loop
	for {
		line, err := rl.Readline()
		if err == io.EOF || strings.TrimSpace(line) == "exit" {
			// Ctrl+D pressed, exit
			// for _, process := range tasks.Processes {
			// 	process.CmdChan <- "stop"
			// }
			break
		}
		if len(line) > 0 {
			cli.CommandHandler(line, tasks)
		}
	}
}
