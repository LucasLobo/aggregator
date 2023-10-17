package inbound

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/lucaslobo/aggregator-cli/internal/common/closer"
	"github.com/lucaslobo/aggregator-cli/internal/common/logs"
	"github.com/lucaslobo/aggregator-cli/internal/core/application"
	"github.com/lucaslobo/aggregator-cli/internal/core/domain"
)

type FileProcessor struct {
	logger logs.Logger
	app    application.Application
}

func NewFileProcessor(logger logs.Logger, app application.Application) FileProcessor {
	return FileProcessor{
		logger: logger,
		app:    app,
	}
}

// CalculateMovingAverageFromFile calculates the moving average with a certain windowSize for the events stored in the
// relative path in filename.
func (f FileProcessor) CalculateMovingAverageFromFile(filename string, windowSize int) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer closer.Close(f.logger, file)

	f.app.Init(windowSize)

	// Let's scan the input file line by line to avoid storing the full file in memory
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Bytes()
		var event domain.TranslationDelivered
		if err = json.Unmarshal(line, &event); err != nil {
			return fmt.Errorf("failed to decode line as JSON: %w", err)
		}
		if err = f.app.ProcessEvent(event); err != nil {
			return fmt.Errorf("error while processing event: %w", err)
		}
	}
	if err = scanner.Err(); err != nil {
		return fmt.Errorf("error scanning file: %w", err)
	}

	return nil
}
