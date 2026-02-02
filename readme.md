# Taskmaster

![42 School](https://img.shields.io/badge/42-Project-000000?style=for-the-badge&logo=42&logoColor=white)
![Go Version](https://img.shields.io/badge/Go-1.22.5-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/github/license/UBA-code/taskmaster?style=for-the-badge&color=green)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen?style=for-the-badge)

A lightweight, production-ready process supervisor and job control daemon written in Go. Taskmaster provides robust process management capabilities with automatic restart policies, graceful shutdowns, and real-time monitoring—all through an intuitive command-line interface.

Built as a modern alternative to [Supervisor](http://supervisord.org/), Taskmaster leverages Go's powerful concurrency primitives to efficiently manage multiple processes without the complexity of traditional thread-based systems.

---

## 📑 Table of Contents

- [Features](#-features)
- [Prerequisites](#-prerequisites)
- [Quick Start](#-quick-start)
- [Configuration](#%EF%B8%8F-configuration)
- [Commands](#-commands)
- [Architecture](#-architecture)
- [Examples](#-examples)
- [Contributing](#-contributing)
- [License](#-license)
- [Authors](#-authors)

---

## ✨ Features

### Core Capabilities

- **🎯 Process Lifecycle Management**: Start, stop, restart, and monitor processes with fine-grained control
- **🔄 Intelligent Restart Policies**: Configure automatic restarts with `always`, `never`, or `on-failure` strategies
- **⚡ Hot Configuration Reload**: Update process definitions without restarting the daemon (via `SIGHUP` or `reload` command)
- **🔍 Real-time Monitoring**: Track process status, PID, uptime, and restart counts
- **📊 Multiple Instances**: Launch multiple instances of the same program with automatic naming
- **📝 Comprehensive Logging**: Separate stdout/stderr logs for each process with tail support
- **🛡️ Graceful Shutdown**: Handle termination signals with configurable timeout and fallback to forced kill
- **🎨 Interactive Shell**: User-friendly CLI with command history powered by readline
- **🔧 Exit Code Handling**: Define expected exit codes to distinguish successful exits from failures
- **🌍 Environment Management**: Set custom environment variables and working directories per process
- **🔐 Permission Control**: Configure umask for process file creation permissions

### Technical Highlights

- **Concurrent Process Monitoring**: Each process runs in its own goroutine for maximum efficiency
- **Channel-based Communication**: Non-blocking message passing between controllers and monitors
- **Signal Propagation**: Forward system signals (SIGTERM, SIGKILL, SIGINT, etc.) to child processes
- **Panic Recovery**: Global panic handler ensures daemon stability
- **Resource Cleanup**: Automatic file descriptor and goroutine cleanup on exit

---

## 🛠 Prerequisites

- **Operating System**: Linux or macOS (UNIX-like systems)
- **Go**: Version 1.18 or higher
- **Dependencies**:
  - `github.com/chzyer/readline` - Interactive shell
  - `gopkg.in/yaml.v3` - YAML configuration parsing

---

## 🚀 Quick Start

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/UBA-code/taskmaster.git
   cd taskmaster
   ```

2. **Build the binary**

   ```bash
   make build
   ```

   This creates the executable at `./bin/taskmaster`

3. **Run with example configuration**
   ```bash
   ./bin/taskmaster config-example.yaml
   ```
   If you run without a config file, Taskmaster will generate `config-example.yaml` for you:
   ```bash
   ./bin/taskmaster
   ```

### Quick Test

After launching Taskmaster, try these commands:

```bash
Taskmaster> status          # View all processes
Taskmaster> start pinger    # Start the example task
Taskmaster> logs pinger 20  # View last 20 log lines
Taskmaster> stop pinger     # Stop the task
Taskmaster> exit            # Gracefully shutdown
```

---

## ⚙️ Configuration

Taskmaster uses YAML configuration files to define process specifications. Each task can have multiple configuration parameters to control its behavior.

### Complete Configuration Example

```yaml
tasks:
  web_server:
    command: "/usr/local/bin/nginx -g 'daemon off;'"
    instances: 1
    autoLaunch: true
    restart: on-failure
    expectedExitCodes: [0]
    successfulStartTimeout: 3
    restartsAttempts: 3
    stopingSignal: SIGTERM
    gracefulStopTimeout: 10
    stdout: /var/log/taskmaster/nginx.out.log
    stderr: /var/log/taskmaster/nginx.err.log
    environment:
      PORT: "8080"
      ENV: "production"
    workingDirectory: /var/www
    unmask: "022"

  worker:
    command: "python3 worker.py"
    instances: 5
    autoLaunch: true
    restart: always
    expectedExitCodes: [0, 2]
    successfulStartTimeout: 1
    restartsAttempts: 5
    stopingSignal: SIGTERM
    gracefulStopTimeout: 15
    stdout: /var/log/taskmaster/worker.out.log
    stderr: /var/log/taskmaster/worker.err.log
    environment:
      WORKER_ID: "auto"
      QUEUE_NAME: "tasks"
    workingDirectory: /opt/app
    unmask: "022"
```

### Configuration Reference

| Parameter                | Type    | Default      | Description                                                        |
| ------------------------ | ------- | ------------ | ------------------------------------------------------------------ |
| `command`                | string  | **required** | Shell command to execute the process                               |
| `instances`              | integer | `1`          | Number of concurrent process instances to launch                   |
| `autoLaunch`             | boolean | `false`      | Start process automatically when Taskmaster launches               |
| `restart`                | string  | `"never"`    | Restart policy: `always`, `never`, or `on-failure`                 |
| `expectedExitCodes`      | int[]   | `[0]`        | Exit codes considered successful (used with `on-failure`)          |
| `successfulStartTimeout` | integer | `5`          | Seconds process must run to be considered "successfully started"   |
| `restartsAttempts`       | integer | `3`          | Maximum restart attempts (0 = unlimited for `always` policy)       |
| `stopingSignal`          | string  | `"SIGTERM"`  | Signal for graceful shutdown: `SIGTERM`, `SIGKILL`, `SIGINT`, etc. |
| `gracefulStopTimeout`    | integer | `10`         | Seconds to wait for graceful stop before forcing SIGKILL           |
| `stdout`                 | string  | `""`         | File path for standard output logs (empty = `/dev/null`)           |
| `stderr`                 | string  | `""`         | File path for standard error logs (empty = `/dev/null`)            |
| `environment`            | map     | `{}`         | Environment variables to set for the process                       |
| `workingDirectory`       | string  | `"."`        | Working directory for process execution                            |
| `unmask`                 | string  | `"022"`      | Umask (permission mask) for files created by process               |

### Configuration Notes

- **Restart Policies**:
  - `always`: Restart regardless of exit code
  - `on-failure`: Restart only if exit code is not in `expectedExitCodes`
  - `never`: Never restart automatically

- **Multiple Instances**: When `instances > 1`, processes are named with suffixes (e.g., `worker_1`, `worker_2`, ...)

- **Hot Reload**: Send `SIGHUP` signal or use `reload` command to apply configuration changes without stopping processes

---

## 💻 Commands

Taskmaster provides an interactive shell with the following commands:

| Command               | Description                                | Example          |
| --------------------- | ------------------------------------------ | ---------------- |
| `status`              | Display status table of all processes      | `status`         |
| `start <task>`        | Start a specific task                      | `start nginx`    |
| `start all`           | Start all configured tasks                 | `start all`      |
| `stop <task>`         | Stop a specific task                       | `stop nginx`     |
| `stop all`            | Stop all running tasks                     | `stop all`       |
| `restart <task>`      | Restart a specific task                    | `restart worker` |
| `restart all`         | Restart all running tasks                  | `restart all`    |
| `reload`              | Reload configuration file                  | `reload`         |
| `logs <task> [lines]` | Display last N lines of logs (default: 10) | `logs worker 50` |
| `help`                | Show command help                          | `help`           |
| `exit`                | Gracefully shutdown all processes and exit | `exit`           |

### Status Display

The `status` command shows a color-coded table with:

- **Task Name**: Process identifier (with instance number for multiple instances)
- **Status**: Current state (see status codes below)
- **PID**: Process ID (or `-` if not running)
- **Uptime**: Duration since start (or `-` if not running)
- **Restarts**: Number of times process has been restarted
- **Command**: The command being executed

#### Status Codes

- 🟢 **RUNNING**: Process running and confirmed started (after `successfulStartTimeout`)
- 🟡 **STARTED**: Process launched but still in startup grace period
- 🟣 **STOPPED**: Process intentionally stopped or not yet started
- 🔴 **FATAL**: Process failed to start or crashed without restart

### Example Session

```bash
$ ./bin/taskmaster config.yaml
Taskmaster> status
Task          Status    PID    Uptime       Restarts  Command
nginx         RUNNING   1234   2m15s        0         /usr/local/bin/nginx
worker_1      RUNNING   1235   2m15s        1         python3 worker.py
worker_2      STOPPED   -      -            0         python3 worker.py

Taskmaster> start worker_2
Process 'worker_2' started with PID 1240

Taskmaster> logs worker_1 5
[2026-02-02 10:15:23] Worker initialized
[2026-02-02 10:15:24] Connected to queue
[2026-02-02 10:15:25] Processing task #1
[2026-02-02 10:15:26] Task completed
[2026-02-02 10:15:27] Waiting for tasks...

Taskmaster> reload
Configuration reloaded.

Taskmaster> stop all
Process 'nginx' stopped gracefully
Process 'worker_1' stopped gracefully
Process 'worker_2' stopped gracefully

Taskmaster> exit
Exiting Taskmaster...
```

---

## 🏗 Architecture

Taskmaster is built on Go's concurrency primitives for efficient, scalable process management.

### Core Components

```
┌─────────────────────────────────────────────────────────────┐
│                        Main Process                          │
│  ┌──────────────┐  ┌─────────────┐  ┌──────────────────┐   │
│  │   CLI Loop   │  │   Config    │  │  Signal Handler  │   │
│  │  (readline)  │  │   Parser    │  │    (SIGHUP)      │   │
│  └──────────────┘  └─────────────┘  └──────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │
                ┌───────────┴───────────┐
                │    Task Manager       │
                │  (Tasks Structure)    │
                └───────────┬───────────┘
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
   ┌────▼────┐         ┌────▼────┐        ┌────▼────┐
   │ Process │         │ Process │        │ Process │
   │ Monitor │         │ Monitor │   ...  │ Monitor │
   │  (Go-   │         │  (Go-   │        │  (Go-   │
   │ routine)│         │ routine)│        │ routine)│
   └────┬────┘         └────┬────┘        └────┬────┘
        │                   │                   │
   ┌────▼────┐         ┌────▼────┐        ┌────▼────┐
   │  Child  │         │  Child  │        │  Child  │
   │ Process │         │ Process │        │ Process │
   └─────────┘         └─────────┘        └─────────┘
```

### Key Design Patterns

#### 1. Goroutine per Process

Each managed process has a dedicated goroutine (`StartTaskManager`) that:

- Listens for commands via buffered channel (`CmdChan`)
- Monitors process lifecycle events
- Handles restart logic independently
- Manages graceful shutdown with timeout

#### 2. Channel-based Communication

- **Command Channels**: Send control signals (start/stop/restart) to process monitors
- **Done Channels**: Receive process exit notifications
- **Timeout Channels**: Implement deadline-based logic (startup success, graceful shutdown)

#### 3. Wait Groups for Synchronization

- **Global WaitGroup** (`tasks.WaitGroup`): Tracks all active processes
- **Process WaitGroup** (`process.Wg`): Tracks instances within a task
- Ensures clean shutdown: all goroutines complete before exit

#### 4. Signal Handling

- **SIGHUP**: Triggers configuration reload without process restart
- **SIGTERM**: Graceful daemon shutdown
- Signals are forwarded to child processes with configurable signal types

### Process State Machine

```
        [STOPPED]
            │
            │ start command
            ▼
        [STARTED] ───────────────────────────┐
            │                                 │
            │ runs > successfulStartTimeout   │ exit with
            ▼                                 │ unexpected code
        [RUNNING]                             │
            │                                 │
            │ stop/restart command            │
            │ or exit                         │
            ▼                                 ▼
        [STOPPED] ◄──────────────────────[FATAL]
                                             │
                        restart policy ──────┘
                        (if applicable)
```

### File Structure

```
taskmaster/
├── cmd/
│   └── taskmaster/
│       └── main.go              # Entry point, CLI loop, signal handlers
├── internal/
│   ├── cli/
│   │   ├── commandHandler.go    # Command parsing and dispatch
│   │   ├── printStatus.go       # Status table formatting
│   │   ├── printLogs.go         # Log file reading
│   │   ├── startProcess.go      # Process lifecycle management
│   │   ├── tasksStruct.go       # Core data structures
│   │   ├── reloadConfig.go      # Hot reload logic
│   │   └── readlineInitializer.go # Shell initialization
│   ├── config/
│   │   ├── configStruct.go      # YAML structure definitions
│   │   ├── configParser.go      # YAML parsing and validation
│   │   └── generateConfig.go    # Example config generation
│   └── logger/
│       └── loggers.go           # Logging utilities
├── go.mod
├── Makefile
└── readme.md
```

---

## 📚 Examples

### Example 1: Simple Web Server

```yaml
tasks:
  api_server:
    command: "node server.js"
    instances: 1
    autoLaunch: true
    restart: always
    successfulStartTimeout: 2
    stdout: /var/log/api.log
    stderr: /var/log/api.err.log
    environment:
      NODE_ENV: "production"
      PORT: "3000"
```

### Example 2: Background Workers

```yaml
tasks:
  email_worker:
    command: "python3 email_worker.py"
    instances: 3
    autoLaunch: true
    restart: on-failure
    expectedExitCodes: [0]
    restartsAttempts: 5
    gracefulStopTimeout: 20
    stdout: /var/log/email_worker.log
    stderr: /var/log/email_worker.err.log
```

### Example 3: Development Environment

```yaml
tasks:
  frontend:
    command: "npm run dev"
    instances: 1
    autoLaunch: false
    restart: never
    stdout: /tmp/frontend.log
    workingDirectory: /home/user/project/frontend

  backend:
    command: "go run main.go"
    instances: 1
    autoLaunch: false
    restart: never
    stdout: /tmp/backend.log
    workingDirectory: /home/user/project/backend
    environment:
      DEBUG: "true"
```

---

## 🤝 Contributing

Contributions are welcome! This project was developed as part of the 42 School curriculum to explore:

- UNIX process management and signals
- Go concurrency patterns (goroutines, channels)
- System programming and daemon design
- I/O multiplexing and non-blocking operations

Feel free to submit issues or pull requests for:

- Bug fixes
- Feature enhancements
- Documentation improvements
- Test coverage

---

## 📄 License

This project is open source. See the LICENSE file for details.

---

## 👥 Authors

Developed with ❤️ by 42 School students:

- **[Yassine Bel Hachmi](https://github.com/UBA-code)** (ybel-hac@student.1337.ma)
- **[Hassan Idhmmououhya](https://github.com/hidhmmou)** (hidhmmou@student.1337.ma)

---

## 🙏 Acknowledgments

- [Supervisor](http://supervisord.org/) - Inspiration for process supervision patterns
- [42 School](https://42.fr/) - Educational framework and project requirements
- Go Community - For excellent concurrency primitives and standard library

---

**⭐ Star this repository if you find it useful!**
