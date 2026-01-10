package cli

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"text/tabwriter"
	"time"
)

func PrintStatus(tasks *Tasks) {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "Task\tStatus\tPID\tUptime\tRestarts\tCommand\n")
	for taskName, process := range tasks.Processes {
		for taskInstanceName, intanceProcess := range process.Instances {
			uptime := formatDuration(intanceProcess.Uptime)
			if !slices.Contains([]string{"RUNNING", "STARTED"}, intanceProcess.Status) {
				uptime = "-"
			}

			// Print task name
			if len(process.Instances) > 1 {
				fmt.Fprintf(w, "%s:%s\t", taskName, taskInstanceName)
			} else {
				fmt.Fprintf(w, "%s\t", taskName)
			}

			// Print status
			fmt.Fprintf(w, "%s\t", intanceProcess.Status)

			// Print PID
			if !slices.Contains([]string{"RUNNING", "STARTED"}, intanceProcess.Status) {
				fmt.Fprintf(w, "-\t")
			} else {
				fmt.Fprintf(w, "%d\t", intanceProcess.Pid)
			}

			// Print uptime, restarts, and command
			fmt.Fprintf(w, "%s\t", uptime)
			fmt.Fprintf(w, "%d\t", intanceProcess.Restarts)
			fmt.Fprintf(w, "%s\n", intanceProcess.Task.Command)
		}
	}

	w.Flush()

	// Apply colors to the formatted output
	lines := strings.Split(buf.String(), "\n")
	for i, line := range lines {
		if i == 0 {
			// Header line in bright white
			fmt.Printf("\033[1;37m%s\033[0m\n", line)
		} else if line != "" {
			// Determine color based on status
			var colorCode string
			if strings.Contains(line, "STARTED") {
				colorCode = "\033[33m" // yellow
			} else if strings.Contains(line, "RUNNING") {
				colorCode = "\033[32m" // green
			} else if strings.Contains(line, "FATAL") {
				colorCode = "\033[31m" // red
			} else if strings.Contains(line, "STOPPED") {
				colorCode = "\033[35m" // magenta
			}

			if colorCode != "" {
				fmt.Printf("%s%s\033[0m\n", colorCode, line)
			} else {
				fmt.Println(line)
			}
		}
	}
}

func formatDuration(d string) string {
	uptimeTime, err := time.Parse(time.RFC3339, d)
	if err != nil {
		return "-"
	}
	dur := time.Since(uptimeTime).Round(time.Second)
	return dur.String()
}
