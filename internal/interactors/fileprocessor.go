package interactors

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

func (f FileProcessor) CalculateMovingAverageFromFile(filename string, windowSize int) error {

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer closer.Close(f.logger, file)

	// We can't scan it in one go because it doesn't have the proper json format
	// Let's scan line by line instead
	scanner := bufio.NewScanner(file)
	var events []domain.TranslationDelivered

	for scanner.Scan() {
		line := scanner.Bytes()
		var event domain.TranslationDelivered
		if err := json.Unmarshal(line, &event); err != nil {
			return fmt.Errorf("failed to decode line as JSON: %w", err)
		}
		events = append(events, event)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning file: %w", err)
	}

	err = f.app.CalculateMovingAverage(events, windowSize)
	if err != nil {
		return fmt.Errorf("failed to calculate moving average: %w", err)
	}
	return nil
}
