package outbound

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/lucaslobo/aggregator/internal/common/logs"
	"github.com/lucaslobo/aggregator/internal/core/domain"
)

// FileWriter is an implementation of a MovingAverageStorer that writes to a file.
type FileWriter struct {
	logger logs.Logger
	folder string

	file           *os.File
	encoder        *json.Encoder
	outputFilePath string
}

func NewFileWriter(logger logs.Logger, folder string) *FileWriter {
	return &FileWriter{
		logger: logger,
		folder: folder,
	}
}

func (f *FileWriter) Close() error {
	file := f.file
	f.file = nil
	f.encoder = nil
	if file != nil {
		return file.Close()
	}
	return nil
}

func (f *FileWriter) setupJSONEncoder() error {
	if f.encoder != nil && f.file != nil {
		return nil
	}

	outputDir, err := createDir(f.folder)
	if err != nil {
		return err
	}

	f.outputFilePath = getOutputPath(outputDir, "events")

	file, err := os.Create(f.outputFilePath)
	if err != nil {
		return err
	}

	f.file = file
	f.encoder = json.NewEncoder(file)
	return nil
}

func (f *FileWriter) StoreMovingAverage(dt domain.AverageDeliveryTime) error {
	err := f.setupJSONEncoder()
	if err != nil {
		return err
	}
	if err = f.encoder.Encode(dt); err != nil {
		return err
	}
	return nil
}

func (f *FileWriter) StoreMovingAverageSlice(deliveryTimes []domain.AverageDeliveryTime) error {
	err := f.setupJSONEncoder()
	if err != nil {
		return err
	}

	for _, dt := range deliveryTimes {
		if err = f.encoder.Encode(dt); err != nil {
			return err
		}
	}

	return nil
}

func createDir(dir string) (string, error) {
	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}
	return dir, nil
}

func getOutputPath(dirpath, filename string) string {
	now := time.Now()
	datetime := now.Format("20060102150405")
	outputFileName := fmt.Sprintf("%s_%s%s", filename, datetime, ".json")
	outputFilePath := fmt.Sprintf("%s/%s", dirpath, outputFileName)
	return outputFilePath
}
