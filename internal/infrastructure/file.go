package infrastructure

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lucaslobo/aggregator-cli/internal/common/closer"
	"github.com/lucaslobo/aggregator-cli/internal/common/logs"
	"github.com/lucaslobo/aggregator-cli/internal/core/domain"
)

// FileWriter is an implementation of a MovingAverageStorer that writes to a file.
type FileWriter struct {
	logger   logs.Logger
	filename string
}

func NewFileWriter(logger logs.Logger, filename string) FileWriter {
	return FileWriter{
		logger:   logger,
		filename: filename,
	}
}

func (f FileWriter) StoreMovingAverage(deliveryTimes []domain.AverageDeliveryTime) error {
	outputDir, err := createDir(f.filename, "output")
	if err != nil {
		return err
	}

	outputFilePath := getOutputPath(f.filename, outputDir)

	file, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer closer.Close(f.logger, file)

	encoder := json.NewEncoder(file)

	f.logger.Infow("Writing to file...",
		"path", outputFilePath)

	// We write to the file line by line because it's not the standard json format
	for _, dt := range deliveryTimes {
		if err := encoder.Encode(dt); err != nil {
			return err
		}
	}

	return nil
}

func createDir(filename, dirname string) (string, error) {
	dir := filepath.Dir(filename)
	outputDir := filepath.Join(dir, dirname)

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", err
	}
	return outputDir, nil
}

func getOutputPath(filename, dirname string) string {
	basename := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	now := time.Now()
	datetime := now.Format("20060102150405")
	outputFileName := fmt.Sprintf("%s_%s%s", basename, datetime, ".json")
	outputFilePath := fmt.Sprintf("%s/%s", dirname, outputFileName)
	return outputFilePath
}
