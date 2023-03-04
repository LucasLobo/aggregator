package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lucaslobo/aggregator-cli/internal/common/logs"
	"github.com/lucaslobo/aggregator-cli/internal/core/application"
	"github.com/lucaslobo/aggregator-cli/internal/core/infrastructureprt"
	"github.com/lucaslobo/aggregator-cli/internal/infrastructure"
	"github.com/lucaslobo/aggregator-cli/internal/interactors"
	"github.com/urfave/cli/v2"
)

const inputFileFlagPropName = "input_file"
const windowSizeFlagPropName = "window_size"
const outputTypeFlagPropName = "output_type"

const fileOutputType = "file"
const stdoutOutputType = "stdout"

// MovingAverageCommand is the command to calculate the moving average aggregation from a file
var MovingAverageCommand = &cli.Command{
	Name:   "moving-average",
	Action: runMovingAverageCommand,
	Flags: []cli.Flag{
		&cli.StringFlag{Name: inputFileFlagPropName, Required: true, Usage: "File (.json) that contains input events"},
		&cli.IntFlag{Name: windowSizeFlagPropName, Required: true, Usage: "Moving average window size in minutes"},
		&cli.StringFlag{Name: outputTypeFlagPropName, Required: false, Usage: fmt.Sprintf("Output type (%s or %s)", fileOutputType, stdoutOutputType), Value: fileOutputType},
	},
}

func runMovingAverageCommand(ctx *cli.Context) error {
	start := time.Now()

	logger, ok := ctx.App.Metadata["Logger"].(logs.Logger)
	if !ok {
		return errors.New("could not get logger")
	}

	inputFile := ctx.String(inputFileFlagPropName)
	windowSize := ctx.Int(windowSizeFlagPropName)
	outputType := ctx.String(outputTypeFlagPropName)

	if windowSize < 1 {
		logger.Warnw("window size cannot be < 1, using default value of 10")
		windowSize = 10
	}

	inputFile = strings.TrimSpace(inputFile)
	if inputFile == "" {
		return errors.New("must provide input file")
	}

	if outputType != fileOutputType && outputType != stdoutOutputType {
		outputType = stdoutOutputType
		logger.Warnf("output type must be file or stdout, using default value of %s", stdoutOutputType)
	}

	logger.Infow("Running Moving Average Command",
		inputFileFlagPropName, inputFile,
		windowSizeFlagPropName, windowSize)

	var storer infrastructureprt.MovingAverageStorer
	if outputType == fileOutputType {
		storer = infrastructure.NewFileWriter(logger, inputFile)
	} else {
		storer = infrastructure.NewStdOut()
	}
	defer storer.Close()
	app := application.New(storer)

	fileProcessor := interactors.NewFileProcessor(logger, app)

	err := fileProcessor.CalculateMovingAverageFromFile(inputFile, windowSize)
	if err != nil {
		return err
	}

	elapsed := time.Since(start)
	logger.Infow("Successfully calculated moving average from file", "time", elapsed)

	return nil
}
