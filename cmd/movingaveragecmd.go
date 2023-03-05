package cmd

import (
	"errors"
	"strings"
	"time"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	awsSqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/urfave/cli/v2"

	"github.com/lucaslobo/aggregator-cli/internal/common/logs"
	"github.com/lucaslobo/aggregator-cli/internal/common/sqs"
	"github.com/lucaslobo/aggregator-cli/internal/core/application"
	"github.com/lucaslobo/aggregator-cli/internal/core/infrastructureprt"
	"github.com/lucaslobo/aggregator-cli/internal/infrastructure"
	"github.com/lucaslobo/aggregator-cli/internal/interactors"
)

const windowSizeFlagPropName = "window_size"
const inputFileFlagPropName = "input_file"
const outputFolderFlagPropName = "output_folder"
const inputQueueFlagPropName = "queue_url"

type cmdCfg struct {
	logger logs.Logger

	windowSize   int
	queueURL     string
	inputFile    string
	outputFolder string

	storer infrastructureprt.MovingAverageStorer
	app    application.Application
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

	defer cfg.storer.Close()

	if cfg.inputFile != "" {
		err = processFromFile(ctx, cfg)
	} else if cfg.queueURL != "" {
		err = processFromQueue(ctx, cfg)
	}

	if err != nil {
		return errors.New("error processing input:" + err.Error())
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
	if outputFolder == "" {
		logger.Warn("Output folder not provided, writing to stdout instead")
	}

	var storer infrastructureprt.MovingAverageStorer
	if outputFolder != "" {
		storer = infrastructure.NewFileWriter(logger, outputFolder)
	} else {
		storer = infrastructure.NewStdOut()
	}

	app := application.New(storer)

	cfg := cmdCfg{
		logger:       logger,
		windowSize:   windowSize,
		inputFile:    inputFile,
		outputFolder: outputFolder,
		queueURL:     queueURL,
		storer:       storer,
		app:          app,
	}

	return cfg, nil
}

func processFromFile(_ *cli.Context, cfg cmdCfg) error {
	cfg.logger.Infow("Running Moving Average Command",
		inputFileFlagPropName, cfg.inputFile,
		windowSizeFlagPropName, cfg.windowSize)

	start := time.Now()
	fileProcessor := interactors.NewFileProcessor(cfg.logger, cfg.app)

	err := fileProcessor.CalculateMovingAverageFromFile(cfg.inputFile, cfg.windowSize)
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
		SqsClient:           *sqsClient,
		SqsURL:              cfg.queueURL,
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     15,
	}
	q := sqs.NewClient(queueCfg)

	queueConsumer := interactors.NewQueueConsumer(cfg.logger, q, cfg.app)
	cfg.logger.Info("Message poller starting...")
	queueConsumer.PollAndProcess(ctx.Context, cfg.windowSize)

	return nil
}
