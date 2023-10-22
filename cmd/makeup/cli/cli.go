package cli

import (
	"fmt"
	"os"
)

type Command func([]string) error

var root Command
var commands = map[string]Command{}

// Setup sets up the CLI with a root command and subcommands
func Setup(rootCmd Command, cmds map[string]Command) {
	root = rootCmd
	commands = cmds
}

// Run runs the CLI with super barebones arg parsing
func Run() error {
	args := os.Args[1:]
	var cmd Command
	var ok bool

	switch len(os.Args) {
	case 1:
		cmd = root
	default:
		cmdName := os.Args[1]
		cmd, ok = commands[cmdName]
		if !ok {
			return fmt.Errorf("not a valid command: %s", cmdName)
		}

		// trim off the command name
		args = args[1:]
	}

	if err := cmd(args); err != nil {
		return err
	}

	return nil
}
