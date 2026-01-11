package cli

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"syscall"
	"time"

	"github.com/uba-code/taskmaster/internal/logger"
)

func parseSignal(signalName string) os.Signal {
	switch signalName {
	case "SIGTERM":
		return syscall.SIGTERM
	case "SIGKILL":
		return syscall.SIGKILL
	case "SIGINT":
		return syscall.SIGINT
	case "SIGWINCH":
		return syscall.SIGWINCH
	default:
		return syscall.SIGTERM
	}
}

func shouldRestart(p *Process, exitCode int) bool {
	isFailure := !slices.Contains(p.Task.ExpectedExitCodes, exitCode) && (p.Status == "STARTED" || p.Status == "FATAL")
	return p.Task.Restart == "always" ||
		(p.Task.Restart == "on-failure" &&
			isFailure &&
			p.Restarts < p.Task.RestartsAttempts)
}

func startStartTimeoutTimer(p *Process) <-chan time.Time {
	successTimeout := make(<-chan time.Time)
	if p.Task.SuccessfulStartTimeout > 0 {
		successTimeout = time.After(time.Duration(p.Task.SuccessfulStartTimeout) * time.Second)
	} else {
		p.Status = "RUNNING"
		logger.Info(fmt.Sprintf("Process '%s' has started successfully", p.Name))
	}
	return successTimeout
}

func startProcess(command string, p *Process, tasks *Tasks) (*exec.Cmd, chan error) {
	cmd := exec.Command("sh", "-c", "umask "+p.Task.Unmask+" && exec "+command)
	env := os.Environ()

	cmd.Stdout = p.Stdout
	cmd.Stderr = p.Stderr

	//? Set environment variables
	for key, value := range p.Task.Environment {
		entry := fmt.Sprintf("%s=%s", key, value)
		env = append(env, entry)
	}
	cmd.Env = env
	cmd.Dir = p.Task.WorkingDirectory

	if err := cmd.Start(); err != nil {
		logger.Info(fmt.Sprintf("'%s' failed to start: %v", p.Name, err))
		p.Status = "FATAL"
		if shouldRestart(p, getExitCode(cmd)) {
			p.Restarts++
			return startProcess(command, p, tasks)
		} else {
			return nil, nil
		}
	}

	tasks.WaitGroup.Add(1)
	p.ParentWg.Add(1)
	p.Status = "STARTED"
	p.Pid = cmd.Process.Pid
	p.Uptime = time.Now().Format(time.RFC3339)
	logger.Success(fmt.Sprintf("Process '%s' started with PID %d", p.Name, p.Pid))

	//? Wait for the process to finish
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	return cmd, done
}

func (p *Process) StartTaskManager(autoStart bool, tasks *Tasks) {
	go func() {
		for {
			cmdReceived := "start"
			if !autoStart {
				cmdReceived = <-p.CmdChan
			}
			autoStart = false
			if (cmdReceived == "stop" || cmdReceived == "restart") && p.Status != "RUNNING" && p.Status != "STARTED" {
				logger.Error(fmt.Sprintf("Process '%s' is not running", p.Name))
			}

			//? if command is not "start", ignore for now
			if cmdReceived != "start" {
				continue
			}

			p.Restarts = 0
			cmd, done := startProcess(p.Task.Command, p, tasks)
			if done == nil {
				continue
			}

			successTimeout := startStartTimeoutTimer(p)

			running := true
			for running {
				select {
				//? if process started successfully after timeout
				case <-successTimeout:
					{
						if p.Status != "RUNNING" {
							p.Status = "RUNNING"
							logger.Info(fmt.Sprintf("Process '%s' has started successfully", p.Name))
						}
					}
				//? if process exits naturally or with error
				case err := <-done:
					{
						tasks.WaitGroup.Done()
						p.ParentWg.Done()
						isFailure := !slices.Contains(p.Task.ExpectedExitCodes, getExitCode(cmd))
						if err != nil {
							logger.Error(fmt.Sprintf("Process '%s' exited with error: %v", p.Name, getExitCode(cmd)))
						} else {
							logger.Info(fmt.Sprintf("Process '%s' exited successfully", p.Name))
						}

						//? restart logic
						if p.Task.Restart == "always" ||
							(p.Task.Restart == "on-failure" &&
								isFailure &&
								p.Restarts < p.Task.RestartsAttempts) {
							for i := p.Restarts; i < p.Task.RestartsAttempts || p.Task.Restart == "always"; i++ {
								p.Restarts++
								successTimeout = startStartTimeoutTimer(p)
								cmd, done = startProcess(p.Task.Command, p, tasks)
								if done == nil {
									continue
								}
								break
							}
							continue
						}
						//? if not restarting, update status
						running = false
						p.Status = "STOPPED"
						if isFailure {
							p.Status = "FATAL"
						}
						break
					}

				//? if command received to stop or restart
				case cmdMsg := <-p.CmdChan:
					{
						if cmdMsg == "start" && (p.Status == "RUNNING" || p.Status == "STARTED") {
							logger.Error(fmt.Sprintf("Process '%s' is already running", p.Name))
							continue
						}
						//? send stopping signal
						if cmdMsg == "stop" || cmdMsg == "restart" {
							cmd.Process.Signal(parseSignal(p.Task.StopingSignal))
						}

						//? wait for graceful shutdown
						select {
						//? It exited gracefully
						case <-done:
							logger.Success(fmt.Sprintf("Process '%s' stopped gracefully", p.Name))

						//? Timeout exceeded, force kill
						case <-time.After(time.Duration(p.Task.GracefulStopTimeout) * time.Second): // "stoptime"
							logger.Info(fmt.Sprintf("Process '%s' did not stop in time, killing...", p.Name))
							cmd.Process.Kill()
						}

						tasks.WaitGroup.Done()
						p.ParentWg.Done()
						p.Status = "STOPPED"
						//? restart if needed
						if cmdMsg == "restart" {
							successTimeout = startStartTimeoutTimer(p)
							cmd, done = startProcess(p.Task.Command, p, tasks)
							if done == nil {
								logger.Error(fmt.Sprintf("Process '%s' failed to restart", p.Name))
								running = false
							} else {
								logger.Info(fmt.Sprintf("Process '%s' restarted", p.Name))
								continue
							}
						}

						logger.Info("Process exited with code " + fmt.Sprint(getExitCode(cmd)))
						//? if restart failed
						isFailure := !slices.Contains(p.Task.ExpectedExitCodes, getExitCode(cmd))
						if isFailure {
							p.Status = "FATAL"
						}
						logger.Info(fmt.Sprintf("Process '%s' has been stopped", p.Name))
						running = false
					}
				}
			}
		}
	}()
}

func getExitCode(cmd *exec.Cmd) int {
	if cmd.ProcessState == nil {
		return -1
	}

	// Check if process was killed by signal on Unix systems
	if status, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
		if status.Signaled() {
			// Process was killed by signal
			// Return 128 + signal number (shell convention)
			return 128 + int(status.Signal())
		}
	}

	return cmd.ProcessState.ExitCode()
}
