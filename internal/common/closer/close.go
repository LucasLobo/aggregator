package closer

import (
	"io"

	"github.com/lucaslobo/aggregator-cli/internal/common/logs"
)

func Close(logger logs.Logger, c io.Closer) {
	err := c.Close()
	if err != nil {
		logger.Errorw("error while closing",
			"error", err)
	}
}
