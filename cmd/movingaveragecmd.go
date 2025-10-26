package cmd

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	awsSqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/urfave/cli/v2"

	"github.com/lucaslobo/aggregator/internal/application"
	"github.com/lucaslobo/aggregator/internal/common/closer"
	"github.com/lucaslobo/aggregator/internal/common/logs"
	"github.com/lucaslobo/aggregator/internal/common/sqs"
	"github.com/lucaslobo/aggregator/internal/inbound"
	"github.com/lucaslobo/aggregator/internal/outbound"
)

const (
	// prop names are used to identify values for the CLI commands
	windowSizeFlagPropName   = "window_size"
	inputFileFlagPropName    = "input_file"
	outputFolderFlagPropName = "output_folder"
	inputQueueFlagPropName   = "queue_url"
)

type cmdCfg struct {
	logger logs.Logger

	windowSize   int
	queueURL     string
	inputFile    string
	outputFolder string

	storerCloser io.Closer
	svc          *application.Application
}

// MovingAverageCommand is the command to calculate the moving average aggregation from a file.
var MovingAverageCommand = &cli.Command{
	Name:   "moving-average",
	Action: runMovingAverageCommand,
	Flags: []cli.Flag{
		&cli.IntFlag{Name: windowSizeFlagPropName, Required: true, Usage: "Moving average window size in minutes"},
		&cli.StringFlag{Name: inputFileFlagPropName, Required: false, Usage: "File (.json) that contains input events"},
		&cli.StringFlag{Name: inputQueueFlagPropName, Required: false, Usage: "SQS Queue URL that contains input events"},
		&cli.StringFlag{Name: outputFolderFlagPropName, Required: false, Usage: "Output folder to write output event files"},
	},
}

func runMovingAverageCommand(ctx *cli.Context) error {
	cfg, err := initCmd(ctx)
	if err != nil {
		return err
	}

	defer closer.Close(cfg.logger, cfg.storerCloser)

	if cfg.inputFile != "" {
		err = processFromFile(ctx, cfg)
	} else if cfg.queueURL != "" {
		err = processFromQueue(ctx, cfg)
	}

	if err != nil {
		return fmt.Errorf("error processing input: %w", err)
	}
	return nil
}

func initCmd(ctx *cli.Context) (cmdCfg, error) {
	logger, ok := ctx.App.Metadata["Logger"].(logs.Logger)
	if !ok {
		return cmdCfg{}, errors.New("could not get logger")
	}

	inputFile := ctx.String(inputFileFlagPropName)
	outputFolder := ctx.String(outputFolderFlagPropName)
	queueURL := ctx.String(inputQueueFlagPropName)
	windowSize := ctx.Int(windowSizeFlagPropName)

	if windowSize < 1 {
		logger.Warnw("window size cannot be < 1, using default value of 10")
		windowSize = 10
	}

	inputFile = strings.TrimSpace(inputFile)
	outputFolder = strings.TrimSpace(outputFolder)
	queueURL = strings.TrimSpace(queueURL)

	if inputFile == "" && queueURL == "" {
		return cmdCfg{}, errors.New("must provide either input file or queue URL")
	}
	if inputFile != "" && queueURL != "" {
		return cmdCfg{}, errors.New("cannot provide both input file and queue URL")
	}

	var storerCloser io.Closer
	var app *application.Application
	if outputFolder != "" {
		storer := outbound.NewFileWriter(logger, outputFolder)
		app = application.New(windowSize, storer)
		storerCloser = storer
	} else {
		logger.Warn("Output folder not provided, writing to stdout instead")
		storer := outbound.NewStdOut()
		app = application.New(windowSize, storer)
		storerCloser = storer
	}

	cfg := cmdCfg{
		logger:       logger,
		windowSize:   windowSize,
		queueURL:     queueURL,
		inputFile:    inputFile,
		outputFolder: outputFolder,
		storerCloser: storerCloser,
		svc:          app,
	}

	return cfg, nil
}

func processFromFile(_ *cli.Context, cfg cmdCfg) error {
	cfg.logger.Infow("Running Moving Average Command from file",
		inputFileFlagPropName, cfg.inputFile,
		windowSizeFlagPropName, cfg.windowSize)

	start := time.Now()
	fileProcessor := inbound.NewFileProcessor(cfg.logger, cfg.svc)

	err := fileProcessor.CalculateMovingAverageFromFile(cfg.inputFile)
	if err != nil {
		return err
	}

	elapsed := time.Since(start)
	cfg.logger.Infow("Successfully calculated moving average from file", "time", elapsed)
	return nil
}

func processFromQueue(ctx *cli.Context, cfg cmdCfg) error {
	cfg.logger.Infow("Running Moving Average Command from SQS Queue",
		inputFileFlagPropName, cfg.queueURL,
		windowSizeFlagPropName, cfg.windowSize)

	awsCfg, err := awsConfig.LoadDefaultConfig(ctx.Context)
	if err != nil {
		return errors.New("could not load AWS default config: " + err.Error())
	}

	sqsClient := awsSqs.NewFromConfig(awsCfg)

	queueCfg := sqs.ConfigSQS{
		Logger:              cfg.logger,
		SqsClient:           sqsClient,
		SqsURL:              cfg.queueURL,
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     15,
	}
	q := sqs.NewClient(queueCfg)

	queueConsumer := inbound.NewQueueConsumer(cfg.logger, q, cfg.svc)
	cfg.logger.Info("Message poller starting...")
	queueConsumer.PollAndProcess(ctx.Context)

	return nil
}
