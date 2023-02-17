package main

import (
	"fmt"
	"os"

	"github.com/lucaslobo/aggregator-cli/cmd"
	"github.com/lucaslobo/aggregator-cli/internal/common/logs"
	"github.com/urfave/cli/v2"
)

const (
	// Name of service
	Name = "service-base-go"
)

var (
	// Version of service
	Version = "1.0.0"
)

// this variable is defined in order to Sync the log at the end of main
var l logs.Logger

func main() {
	app := cli.App{
		Name:        Name,
		Version:     Version,
		Description: "Simple command line application that parses a stream of events and produces an aggregated output.",
		Before:      setupBefore,
		Commands: []*cli.Command{
			cmd.MovingAverageCommand,
		},
		DefaultCommand: cmd.MovingAverageCommand.Name,
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)

	err := app.Run(os.Args)
	defer l.Sync() // flushes buffer, if any

	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}

func setupBefore(cli *cli.Context) error {
	logger, err := logs.New()
	if err != nil {
		return err
	}
	cli.App.Metadata["Logger"] = logger

	l = logger
	return nil
}
