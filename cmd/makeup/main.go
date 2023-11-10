package main

import (
	"log/slog"
	"os"

	"github.com/cohix/makeup/cmd/makeup/cli"
	"github.com/cohix/makeup/cmd/makeup/commands"
)

func main() {
	cli.Setup(
		commands.Root,
		map[string]cli.Command{
			"add":   commands.Add,
			"build": commands.Build,
			"test":  commands.Test,
			"clean": commands.Clean,
		},
	)

	if err := cli.Run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
