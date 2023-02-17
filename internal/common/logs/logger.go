package logs

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

func New() (Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return Logger{}, fmt.Errorf("error creating logger: %w", err)
	}

	sugar := logger.Sugar()
	return Logger{
		sugar,
	}, nil
}

func (l Logger) Sync() error {
	if l.SugaredLogger != nil {
		return l.SugaredLogger.Sync()
	}
	return nil
}
