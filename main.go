package main

import (
	"log"

	"github.com/cohix/makeup/cli"
	"github.com/cohix/makeup/commands"
)

func main() {
	cli.Setup(
		commands.Root,
		map[string]cli.Command{
			"test":  commands.Test,
			"clean": commands.Clean,
		},
	)

	if err := cli.Run(); err != nil {
		log.Fatal(err)
	}
}
