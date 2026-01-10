package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/uba-code/taskmaster/internal/logger"
)

func CommandHandler(command string, tasks *Tasks, configFileName string) {
	commandParts := strings.Split(strings.TrimSpace(command), " ")
	command = commandParts[0]

	switch command {
	case "status":
		PrintStatus(tasks)

	case "reload":
		logger.Info("Configuration reloaded.")
		ReloadConfig(tasks, configFileName)

	case "start", "stop", "restart":
		if len(commandParts) < 2 {
			logger.Error("Please specify the task to " + command)
			return
		}
		argument := commandParts[1]
		if argument == "all" {
			if command == "start" { // bonus
				StartAllProcesses(tasks)
			} else if command == "stop" { // bonus
				StopAllProcesses(tasks)
			} else if command == "restart" { // bonus
				RestartAllProcesses(tasks)
			}
			return
		}
		if process, exists := tasks.Processes[argument]; exists {
			for _, p := range process.Instances {
				p.CmdChan <- command
			}
		} else {
			logger.Error("Task '" + argument + "' not found.")
		}

	case "logs":
		count := 10 // default log lines
		if len(commandParts) > 3 || len(commandParts) < 2 {
			logger.Error("Usage: logs <task-name>")
			return
		} else if len(commandParts) == 3 {
			c, err := strconv.Atoi(commandParts[2])
			if err != nil || c < 1 {
				logger.Error("Please provide a valid number of lines to display.")
				return
			}
			count = c
		}
		printLogs(commandParts[1], count, tasks)
	case "help": // bonus
		fmt.Println("Available commands:")
		fmt.Println("\tstatus\t\t\t\t- Show system status")
		fmt.Println("\treload\t\t\t\t- Reload configuration")
		fmt.Println("\tstart\t\t\t\t- Start the service")
		fmt.Println("\tstart all\t\t\t- Start all services")
		fmt.Println("\trestart\t\t\t\t- Restart the service")
		fmt.Println("\trestart all\t\t\t- Restart all services")
		fmt.Println("\tstop\t\t\t\t- Stop the service")
		fmt.Println("\tstop all\t\t\t- Stop all services")
		fmt.Println("\tlogs [task-name] [lines]\t- Show the last 'lines' of logs for the specified task (default 10)")
		fmt.Println("\texit\t\t\t\t- Exit the application")
		fmt.Println("\thelp\t\t\t\t- Show this help message")
	default:
		fmt.Printf("Unknown command: %s. Type 'help' for a list of commands.\n", command)
	}
}
