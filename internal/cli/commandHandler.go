package cli

import (
	"fmt"
	"strings"

	"github.com/uba-code/taskmaster/internal/logger"
)

func CommandHandler(command string, tasks *Tasks, filename string) {
	commandParts := strings.Split(strings.TrimSpace(command), " ")
	command = commandParts[0]

	switch command {
	case "status":
		PrintStatus(tasks)

	case "reload":
		logger.Info("Configuration reloaded.")
		ReloadConfig(tasks, filename)
	case "start", "stop", "restart":
		if len(commandParts) < 2 {
			logger.Error("Please specify the task to " + command)
			return
		}
		argument := commandParts[1]
		if process, exists := tasks.Processes[argument]; exists {
			for _, p := range process.Instances {
				p.CmdChan <- command
			}
		} else {
			logger.Error("Task '" + argument + "' not found.")
		}

	case "help":
		fmt.Println("Available commands:")
		fmt.Println("\tstatus  - Show system status")
		fmt.Println("\treload  - Reload configuration")
		fmt.Println("\tstart   - Start the service")
		fmt.Println("\trestart - Restart the service")
		fmt.Println("\tstop    - Stop the service")
		fmt.Println("\texit    - Exit the application")
		fmt.Println("\thelp    - Show this help message")
	default:
		fmt.Printf("Unknown command: %s. Type 'help' for a list of commands.\n", command)
	}
}
