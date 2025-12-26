package cli

import (
	"fmt"
	"os"
	"slices"
	"text/tabwriter"
	"time"
)

func PrintStatus(tasks *Tasks) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintf(w, "Task\tStatus\tPID\tUptime\tRestarts\n")
	for taskName, process := range tasks.Processes {
		for taskInstanceName, intanceProcess := range process.Instances {
			uptime := formatDuration(intanceProcess.Uptime)
			if !slices.Contains([]string{"RUNNING", "STARTED"}, intanceProcess.Status) {
				uptime = "-"
			}
			if len(process.Instances) > 1 {
				fmt.Fprintf(w, "%s:%s\t", taskName, taskInstanceName)
			} else {
				fmt.Fprintf(w, "%s\t", taskName)
			}
			fmt.Fprintf(w, "%s\t", intanceProcess.Status)
			if !slices.Contains([]string{"RUNNING", "STARTED"}, intanceProcess.Status) {
				fmt.Fprintf(w, "-\t")
			} else {
				fmt.Fprintf(w, "%d\t", intanceProcess.Pid)
			}
			fmt.Fprintf(w, "%s\t", uptime)
			fmt.Fprintf(w, "%d\n", intanceProcess.Restarts)
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
