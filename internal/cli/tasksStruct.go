package cli

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/uba-code/taskmaster/internal/config"
	"github.com/uba-code/taskmaster/internal/logger"
)

type Tasks struct {
	Processes map[string]*ProcessInstances
	WaitGroup sync.WaitGroup
}

type ProcessInstances struct {
	TaskCfg   *config.TaskCfg
	Instances map[string]*Process
	Wg        sync.WaitGroup
}

type Process struct {
	Status     string // STARTED, RUNNING, STOPPED, FATAL, STARTING, STOPPING
	Pid        int
	Uptime     string
	Name       string
	Stdout     *os.File
	Stderr     *os.File
	Restarts   int
	Task       *config.TaskCfg
	CmdChan    chan string
	StatusChan chan string
	ParentWg   *sync.WaitGroup
}

func NewTasksObj(config *config.Config) *Tasks {
	var taskArr *Tasks = &Tasks{
		Processes: make(map[string]*ProcessInstances),
	}

	for taskName, task := range config.Tasks {
		taskArr.Processes[taskName] = initTask(&task, taskName)
		launchProcessInstances(taskArr.Processes[taskName], taskArr)
	}

	return taskArr
}

func initTask(task *config.TaskCfg, taskName string) *ProcessInstances {
	var process = &ProcessInstances{
		Instances: make(map[string]*Process),
	}
	process.TaskCfg = task
	for i := 0; i < task.Instances; i++ {
		stdout, stderr := setOutputFiles(task.Stdout, task.Stderr)
		processName := taskName
		if task.Instances > 1 {
			processName = processName + "_" + strconv.Itoa(i+1)
		}
		process.Instances[processName] = &Process{
			Status:     "STOPPED",
			Pid:        0,
			Uptime:     "",
			Name:       processName,
			Stdout:     stdout,
			Stderr:     stderr,
			Restarts:   0,
			Task:       task,
			CmdChan:    make(chan string),
			StatusChan: make(chan string),
			ParentWg:   &process.Wg,
		}
	}
	return process
}

func launchProcessInstances(process *ProcessInstances, taskArr *Tasks) {
	for InstanceKey, p := range process.Instances {
		process.Instances[InstanceKey].StartTaskManager(p.Task.AutoLaunch, taskArr)
	}
}

func (p *ProcessInstances) StopAllInstances() {
	for _, instance := range p.Instances {
		instance.CmdChan <- "stop"
	}
	p.Wg.Wait()
}

func setOutputFiles(stdoutPath string, stderrPath string) (stdout *os.File, stderr *os.File) {
	if stdoutPath == "" {
		stdoutPath = "/dev/null"
	}
	if stderrPath == "" {
		stderrPath = "/dev/null"
	}
	stdout, err := os.OpenFile(stdoutPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to open stdout log file: %s", err.Error()))
		stderr = nil
	}
	stderr, err = os.OpenFile(stderrPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to open stderr log file: %s", err.Error()))
		stderr = nil
	}
	return stdout, stderr
}
