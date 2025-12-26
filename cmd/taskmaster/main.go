package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

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
	logger.InitializeLogFile()

	var rl = cli.ReadlineInit()
	defer rl.Close()
	logger.SetReadline(rl)
	var tasks = cli.NewTasksObj(cfg)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGTERM)

	go func() {
		<-signalChan
		cli.ReloadConfig(tasks, os.Args[1])
	}()

	//* cli loop
	for {
		line, err := rl.Readline()
		if err == io.EOF || strings.TrimSpace(line) == "exit" {
			logger.Info("Exiting Taskmaster...")
			for processName, proc := range tasks.Processes {
				for instanceName, p := range proc.Instances {
					if p.Status == "STARTED" || p.Status == "RUNNING" {
						p.CmdChan <- "stop"
					}
					if p.Stderr != nil {
						if err := p.Stderr.Close(); err != nil {
							logger.Error(fmt.Sprintf("Failed to close stderr for process '%s' of task '%s': %v", instanceName, processName, err))
						}
					}
					if p.Stdout != nil {
						if err := p.Stdout.Close(); err != nil {
							logger.Error(fmt.Sprintf("Failed to close stdout for process '%s' of task '%s': %v", instanceName, processName, err))
						}
					}
					logger.Info("Stopping process '" + instanceName + "' of task '" + processName + "'")
				}
			}
			tasks.WaitGroup.Wait()
			break
		}
		if len(line) > 0 {
			cli.CommandHandler(line, tasks, os.Args[1])
		}
	}
}
