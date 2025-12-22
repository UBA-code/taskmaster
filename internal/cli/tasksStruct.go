package cli

import (
	"github.com/uba-code/taskmaster/internal/config"
)

type Tasks struct {
	Processes map[string]*Process
}

type Process struct {
	Status     string // RUNNING, STOPPED, FATAL, STARTING, STOPPING
	Pid        int
	Uptime     string
	Name       string
	Restarts   int
	Task       config.Task
	CmdChan    chan string
	StatusChan chan string
}

func NewTasksObj(config config.Config) *Tasks {
	var taskArr *Tasks = &Tasks{
		Processes: make(map[string]*Process),
	}

	for taskName, task := range config.Tasks {
		taskArr.Processes[taskName] = &Process{
			Status:     "STOPPED",
			Pid:        0,
			Uptime:     "2025-12-21T16:49:10+01:00",
			Task:       task,
			Name:       taskName,
			Restarts:   0,
			CmdChan:    make(chan string),
			StatusChan: make(chan string),
		}
	}

	for key := range taskArr.Processes {
		taskArr.Processes[key].StartTaskManager()
	}

	return taskArr
}
