package config

type TaskCfg struct {
	Command                string            `yaml:"command"`
	Instances              int               `yaml:"instances"`
	AutoLaunch             bool              `yaml:"autoLaunch"`
	Restart                string            `yaml:"restart"`
	ExpectedExitCodes      []int             `yaml:"expectedExitCodes"`
	SuccessfulStartTimeout int               `yaml:"successfulStartTimeout"` // seconds
	RestartsAttempts       int               `yaml:"restartsAttempts"`
	StopingSignal          string            `yaml:"stopingSignal"`
	GracefulStopTimeout    int               `yaml:"gracefulStopTimeout"` // seconds
	Stdout                 string            `yaml:"stdout"`              // pointer to handle null
	Stderr                 string            `yaml:"stderr"`              // pointer to handle null
	Environment            map[string]string `yaml:"environment"`
	WorkingDirectory       string            `yaml:"workingDirectory"`
	Unmask                 string            `yaml:"unmask"`
}

type Config struct {
	Tasks map[string]TaskCfg `yaml:"tasks"`
}

// NewTaskCfg creates a new TaskCfg struct with default values
func NewTaskCfg() TaskCfg {
	return TaskCfg{
		Command:                "",
		Instances:              1,
		AutoLaunch:             false,
		Restart:                "never",
		ExpectedExitCodes:      []int{0},
		SuccessfulStartTimeout: 5,
		RestartsAttempts:       3,
		StopingSignal:          "TERM",
		GracefulStopTimeout:    10,
		Stdout:                 "",
		Stderr:                 "",
		Environment:            make(map[string]string),
		WorkingDirectory:       ".",
		Unmask:                 "022",
	}
}

// SetDefaults sets default values for fields that are not set (zero values)
func (t *TaskCfg) SetDefaults() {
	if t.Instances == 0 {
		t.Instances = 1
	}
	if t.Restart == "" {
		t.Restart = "never"
	}
	if t.ExpectedExitCodes == nil {
		t.ExpectedExitCodes = []int{0}
	}
	if t.SuccessfulStartTimeout == 0 {
		t.SuccessfulStartTimeout = 5
	}
	// if t.RestartsAttempts == 0 {
	// 	t.RestartsAttempts = 0
	// }
	if t.StopingSignal == "" {
		t.StopingSignal = "SIGTERM"
	}
	// if t.GracefulStopTimeout == 0 {
	// 	t.GracefulStopTimeout = 0
	// }
	if t.Environment == nil {
		t.Environment = make(map[string]string)
	}
	if t.WorkingDirectory == "" {
		t.WorkingDirectory = "."
	}
	if t.Unmask == "" {
		t.Unmask = "022"
	}
	// Command, Stdout, Stderr left as is if empty
}
