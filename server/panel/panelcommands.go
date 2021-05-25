package panel

import "fmt"

func (cp *ControlPanel) registerCommands() {
	cp.commands = make(map[string]func([]string))

	// Command: register
	cp.HandleCommand("register", func(args []string) {
		if len(args) != 2 {
			fmt.Println("Usage: register <username> <password>")
			return
		}
			username := args[0]
			password := args[1]
			cp.ChatServer.UserManager.RegisterUser(username,password)
			fmt.Printf("User: %s registered with password: %s\n", username, password)

	})

	// Command: unregister
	cp.HandleCommand("unregister", func(args []string) {
		if len(args) == 0 {
			fmt.Println("Usage: unregister <username>")
			return
		}

		if cp.ChatServer.UserManager.UnregisterUser(args[0]) {
			cp.ChatServer.KickUser(args[0])
			fmt.Printf("The user: %s has been unregistered successfully!\n", args[0])
		} else {
			fmt.Printf("Could not find user: %s\n", args[0])
		}
	})

	cp.HandleCommand("kick", func(args []string) {
		if len(args) == 0 {
			fmt.Println("Usage: kick <username>")
			return
		}

		if cp.ChatServer.KickUser(args[0]) {
			fmt.Printf("The user: %s has been kicked successfully!\n", args[0])
		} else {
			fmt.Printf("Could not find user: %s\n", args[0])
		}
	})
}