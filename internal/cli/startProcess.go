package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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
	default:
		return syscall.SIGTERM
	}
}

func startProcess(command string, p *Process) (*exec.Cmd, chan error) {
	commandParts := strings.Split(command, " ")
	cmd := exec.Command(commandParts[0], commandParts[1:]...)
	if err := cmd.Start(); err != nil {
		logger.Info(fmt.Sprintf("'%s' failed to start: %v", p.Name, err))
		return nil, nil
	}

	p.Status = "RUNNING"
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

func (p *Process) StartTaskManager() {
	go func() {
		for {
			cmdReceived := <-p.CmdChan
			if (cmdReceived == "stop" || cmdReceived == "restart") && p.Status != "RUNNING" {
				logger.Error(fmt.Sprintf("Process '%s' is not running", p.Name))
			}

			//? if command is not "start", ignore for now
			if cmdReceived != "start" {
				continue
			}

			cmd, done := startProcess(p.Task.Command, p)
			if done == nil {
				continue
			}

			running := true
			for running {
				select {
				//? if process exits naturally
				case err := <-done:
					if err != nil {
						logger.Error(fmt.Sprintf("Process '%s' exited with error: %v", p.Name, err))
					} else {
						logger.Info(fmt.Sprintf("Process '%s' exited successfully", p.Name))
					}
					p.Status = "STOPPED"
				//? if command received to stop or restart
				case cmdMsg := <-p.CmdChan:
					if cmdMsg == "start" && p.Status == "RUNNING" {
						logger.Error(fmt.Sprintf("Process '%s' is already running", p.Name))
						continue
					}
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
						cmd.Process.Kill()
					}
					p.Status = "STOPPED"
					if cmdMsg == "restart" {
						cmd, done = startProcess(p.Task.Command, p)
						if done == nil {
							logger.Error(fmt.Sprintf("Process '%s' failed to restart", p.Name))
							running = false
							break
						}
						logger.Info(fmt.Sprintf("Process '%s' restarted", p.Name))
						continue
					}
					logger.Info(fmt.Sprintf("Process '%s' has been stopped", p.Name))
					running = false
				}
			}
		}
	}()
}
