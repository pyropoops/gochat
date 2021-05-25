package panel

import (
	"bufio"
	"fmt"
	"gochat/server/chat"
	"os"
	"strings"
)

type ControlPanel struct {
	ChatServer chat.Server
	commands   map[string]func([]string)
}

func (cp *ControlPanel) HandleCommand(command string, callback func([]string)) {
	cp.commands[command] = callback
}

func (cp *ControlPanel) Start() {
	cp.registerCommands()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		}
		arr := strings.Split(scanner.Text(), " ")
		if command, has := cp.commands[arr[0]]; has {
			command(arr[1:])
		} else {
			fmt.Printf("That command: %s was not found.\n", arr[0])
		}
	}
}

func NewPanel(server chat.Server) ControlPanel {
	return ControlPanel{
		ChatServer: server,
		commands:   make(map[string]func([]string)),
	}
}