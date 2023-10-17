package inbound

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/lucaslobo/aggregator-cli/internal/common/closer"
	"github.com/lucaslobo/aggregator-cli/internal/common/logs"
	"github.com/lucaslobo/aggregator-cli/internal/core/domain"
	"github.com/lucaslobo/aggregator-cli/internal/core/inboundprt"
)

type FileProcessor struct {
	logger logs.Logger
	svc    inboundprt.MovingAverageCalculator
}

func NewFileProcessor(logger logs.Logger, svc inboundprt.MovingAverageCalculator) FileProcessor {
	return FileProcessor{
		logger: logger,
		svc:    svc,
	}
}

// CalculateMovingAverageFromFile calculates the moving average for the events stored in the file (relative path).
func (f FileProcessor) CalculateMovingAverageFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer closer.Close(f.logger, file)

	// Let's scan the input file line by line to avoid storing the full file in memory
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Bytes()
		var event domain.TranslationDelivered
		if err = json.Unmarshal(line, &event); err != nil {
			return fmt.Errorf("failed to decode line as JSON: %w", err)
		}
		if err = f.svc.ProcessEvent(event); err != nil {
			return fmt.Errorf("error while processing event: %w", err)
		}
	}
	if err = scanner.Err(); err != nil {
		return fmt.Errorf("error scanning file: %w", err)
	}

	return nil
}
