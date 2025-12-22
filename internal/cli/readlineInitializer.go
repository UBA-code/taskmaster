package cli

import "github.com/chzyer/readline"

var completer = readline.NewPrefixCompleter(
	readline.PcItem("status"),
	readline.PcItem("reload"),
	readline.PcItem("start"),
	readline.PcItem("restart"),
	readline.PcItem("stop"),
	readline.PcItem("exit"),
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
