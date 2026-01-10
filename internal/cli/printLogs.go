package cli

import (
	"os"

	"github.com/uba-code/taskmaster/internal/logger"
)

func printLogs(taskname string, count int, tasks *Tasks) {
	if process, exists := tasks.Processes[taskname]; exists {
		if process.TaskCfg.Stdout == "" {
			logger.Error("No logs available for this task.")
			return
		}
		data, err := os.ReadFile(process.TaskCfg.Stdout)
		if err != nil {
			logger.Error("Error reading log file: " + err.Error())
			return
		}
		logLines := SplitLines(string(data))
		start := len(logLines) - count
		if start < 0 {
			start = 0
		}
		for _, line := range logLines[start:] {
			println(line)
		}
	}
}

func SplitLines(s string) []string {
	var lines []string
	currentLine := ""
	for _, char := range s {
		if char == '\n' {
			lines = append(lines, currentLine)
			currentLine = ""
		} else {
			currentLine += string(char)
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	return lines
}
