package closer

import (
	"io"

	"github.com/lucaslobo/aggregator-cli/internal/common/logs"
)

// Close the provided io.Closer and log on error
func Close(logger logs.Logger, c io.Closer) {
	err := c.Close()
	if err != nil {
		logger.Errorw("error while closing",
			"error", err)
	}
}
