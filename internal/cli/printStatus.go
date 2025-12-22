package cli

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"
)

func PrintStatus(tasks *Tasks) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintf(w, "Task\tStatus\tPID\tUptime\n")
	for taskName, process := range tasks.Processes {
		uptime := formatDuration(process.Uptime)
		if process.Status != "RUNNING" {
			uptime = "-"
		}
		// fmt.Fprintf(w, "%s\t%s\t%d\t%s\n", taskName, process.Status, pid, uptime)
		fmt.Fprintf(w, "%s\t", taskName)
		fmt.Fprintf(w, "%s\t", process.Status)
		if process.Status != "RUNNING" {
			fmt.Fprintf(w, "-\t")
		} else {
			fmt.Fprintf(w, "%d\t", process.Pid)
		}
		fmt.Fprintf(w, "%s\n", uptime)
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
