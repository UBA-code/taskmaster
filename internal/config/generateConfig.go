package config

import (
	"os"
)

func GenerateConfig() {
	content := `tasks:
  pinger:
    command: "ping 8.8.8.8.8"
    instances: 5
    autoLaunch: true
    restart: on-failure # always, never, on-failure
    expectedExitCodes:
      - 0
    successfulStartTimeout: 1 # seconds
    restartsAttempts: 2
    stopingSignal: SIGTERM
    gracefulStopTimeout: 15 # seconds
    stdout: /var/log/sleeper_stdout.log
    stderr: /var/log/sleeper_stderr.log
    environment:
      ENV_VAR1: value1
      ENV_VAR2: value2
    workingDirectory: /tmp
    unmask: 777
`
	err := os.WriteFile("config-example.yaml", []byte(content), 0644)
	if err != nil {
		panic(err)
	}
}
