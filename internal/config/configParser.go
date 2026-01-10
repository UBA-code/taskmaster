package config

import (
	"os"

	"github.com/uba-code/taskmaster/internal/logger"
	"gopkg.in/yaml.v3"
)

func ParseConfig(filename string) *Config {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic("Failed to read config file: " + err.Error())
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic("Failed to parse config file: " + err.Error())
	}

	for taskName, task := range config.Tasks {
		if taskName == "all" {
			logger.Error("Task name 'all' is reserved and cannot be used.")
			os.Exit(1)
		}
		if task.Command == "" {
			logger.Error("Task '" + taskName + "' is missing a command.")
			os.Exit(1)
		}
	}

	// Set defaults for tasks that have unset fields
	for i := range config.Tasks {
		task := config.Tasks[i]
		task.SetDefaults()
		config.Tasks[i] = task
	}

	return &config
}
