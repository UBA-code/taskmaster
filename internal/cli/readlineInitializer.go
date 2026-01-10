package cli

import "github.com/chzyer/readline"

var completer = readline.NewPrefixCompleter(
	readline.PcItem("status"),
	readline.PcItem("reload"),
	readline.PcItem("start"),
	readline.PcItem("start all"),
	readline.PcItem("restart"),
	readline.PcItem("restart all"),
	readline.PcItem("stop"),
	readline.PcItem("stop all"),
	readline.PcItem("exit"),
	readline.PcItem("logs"),
	readline.PcItem("help"),
)

func ReadlineInit() *readline.Instance {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "TaskMaster> ",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	return rl
}
