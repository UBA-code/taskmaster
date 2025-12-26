package config

import (
	"os"

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

	// Set defaults for tasks that have unset fields
	for _, task := range config.Tasks {
		task.SetDefaults()
	}

	return &config
}
