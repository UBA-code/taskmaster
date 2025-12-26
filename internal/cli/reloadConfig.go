package cli

import (
	"fmt"

	"github.com/uba-code/taskmaster/internal/config"
)

func ReloadConfig(existingTasks *Tasks, filename string) {
	cfg := config.ParseConfig(filename)

	for taskName, task := range cfg.Tasks {
		//? If task exists and is different than old task, update it
		if _, exists := existingTasks.Processes[taskName]; exists && isDifferent(existingTasks.Processes[taskName].TaskCfg, &task) {
			fmt.Printf("Found differences in task '%s', reloading...\n", taskName)
			existingTasks.Processes[taskName].StopAllInstances() //! this will wait for all instances to stop
			existingTasks.Processes[taskName] = &ProcessInstances{
				TaskCfg:   &task,
				Instances: make(map[string]*Process),
			}
			existingTasks.Processes[taskName] = initTask(&task, taskName)
			launchProcessInstances(existingTasks.Processes[taskName], existingTasks)
		} else if !exists {
			existingTasks.Processes[taskName] = &ProcessInstances{
				TaskCfg:   &task,
				Instances: make(map[string]*Process),
			}
			existingTasks.Processes[taskName] = initTask(&task, taskName)
			launchProcessInstances(existingTasks.Processes[taskName], existingTasks)
		}
	}
	//? If task exists in old config but not in new one, stop it
	for taskName, processInstances := range existingTasks.Processes {
		if _, exists := cfg.Tasks[taskName]; !exists {
			fmt.Printf("Task '%s' not found in new config, stopping...\n", taskName)
			processInstances.StopAllInstances() //! this will wait for all instances to stop
			delete(existingTasks.Processes, taskName)
		}
	}
}

func isDifferent(oldTask *config.TaskCfg, newTask *config.TaskCfg) bool {
	diffEnv := len(oldTask.Environment) != len(newTask.Environment)
	for k, v := range oldTask.Environment {
		if val, exists := newTask.Environment[k]; !exists || val != v {
			diffEnv = true
			break
		}
	}
	return oldTask.Command != newTask.Command ||
		oldTask.Instances != newTask.Instances ||
		oldTask.AutoLaunch != newTask.AutoLaunch ||
		oldTask.Restart != newTask.Restart ||
		oldTask.SuccessfulStartTimeout != newTask.SuccessfulStartTimeout ||
		oldTask.RestartsAttempts != newTask.RestartsAttempts ||
		oldTask.StopingSignal != newTask.StopingSignal ||
		oldTask.GracefulStopTimeout != newTask.GracefulStopTimeout ||
		oldTask.Stdout != newTask.Stdout ||
		oldTask.Stderr != newTask.Stderr ||
		diffEnv ||
		oldTask.WorkingDirectory != newTask.WorkingDirectory ||
		oldTask.Unmask != newTask.Unmask
}
