package logs

import (
	"fmt"

	"go.uber.org/zap"
)

// Logger is a simple wrapper around *zap.SugaredLogger
// currently this wrapping doesn't offer many benefits, but it could be used to abstract
// the underlying logger into a common logging interface
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
