package main

import (
	"fmt"
	"os"

	"github.com/ice-stuff/ice-agent/commands"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ice-agent"
	app.Usage = "iCE Agent program."
	app.Version = "2.1.0"

	app.Commands = []cli.Command{
		commands.RegisterSelfCommand,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
}
