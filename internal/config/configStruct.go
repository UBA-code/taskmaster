package config

type Task struct {
	Command                 string            `yaml:"command"`
	Instances               int               `yaml:"instances"`
	AutoLaunch              bool              `yaml:"autoLaunch"`
	Restart                 string            `yaml:"restart"`
	FailureExitCodes        []int             `yaml:"failureExitCodes"`
	SuccessfullStartTimeout int               `yaml:"successfullStartTimeout"` // seconds
	RestartsAttempts        int               `yaml:"restartsAttempts"`
	StopingSignal           string            `yaml:"stopingSignal"`
	GracefulStopTimeout     int               `yaml:"gracefulStopTimeout"` // seconds
	Stdout                  string            `yaml:"stdout"`              // pointer to handle null
	Stderr                  string            `yaml:"stderr"`              // pointer to handle null
	Environment             map[string]string `yaml:"environment"`
	WorkingDirectory        string            `yaml:"workingDirectory"`
	Unmask                  int               `yaml:"unmask"`
}

type Config struct {
	Tasks map[string]Task `yaml:"tasks"`
}
