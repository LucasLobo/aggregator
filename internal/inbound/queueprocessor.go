package inbound

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	awsSQSTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"

	"github.com/lucaslobo/aggregator/internal/common/logs"
	"github.com/lucaslobo/aggregator/internal/domain"
)

type Queue interface {
	GetMessages(ctx context.Context) (*sqs.ReceiveMessageOutput, error)
	SendMessage(ctx context.Context, message awsSQSTypes.Message) error
	Delete(ctx context.Context, message awsSQSTypes.Message) error
	ChangeMessageVisibility(ctx context.Context, receiptHandle *string, timeout int64) error
}

type QueueConsumer struct {
	logger logs.Logger

	queueClient Queue
	svc         MovingAverageCalculator
}

func NewQueueConsumer(logger logs.Logger, queueClient Queue, svc MovingAverageCalculator) QueueConsumer {
	return QueueConsumer{
		logger:      logger,
		queueClient: queueClient,
		svc:         svc,
	}
}

// PollAndProcess polls the queue and processes the messages
func (c *QueueConsumer) PollAndProcess(ctx context.Context) {
	// indefinitely poll queue for messages
	for {
		messages, err := c.readQueueMessages(ctx)
		if errors.Is(err, errNoMessages) {
			c.logger.Info("no messages found")
			continue
		} else if err != nil {
			c.logger.Errorw("unexpected error when reading queue",
				"error", err)
			continue
		}
		c.processMessages(ctx, messages)
	}
}

var errNoMessages = errors.New("no sqs messages found")

func (c *QueueConsumer) readQueueMessages(ctx context.Context) ([]awsSQSTypes.Message, error) {
	res, err := c.queueClient.GetMessages(ctx)
	if err != nil {
		return nil, err
	}
	if res == nil || len(res.Messages) == 0 {
		return nil, errNoMessages
	}

	return res.Messages, nil
}

func (c *QueueConsumer) processMessages(ctx context.Context, messages []awsSQSTypes.Message) {
	c.logger.Infow("read messages from queue", "quantity", len(messages))
	for _, message := range messages {
		if message.Body != nil {
			var event domain.DurationEvent
			err := json.Unmarshal([]byte(*message.Body), &event)
			if err != nil {
				// for simplicity purposes, when an error occurs we simply log it...
				// ideally we'd have a better error handler, such as reporting the error to Sentry/NR
				// and managing the message in the DLQ
				c.logger.Errorw("error unmarshalling message", "error", err)
				return
			}

			err = c.svc.ProcessEvent(event)
			if err != nil {
				c.logger.Errorw("could not process message", "error", err)
				return
			}

			err = c.queueClient.Delete(ctx, message)
			if err != nil {
				c.logger.Errorw("could not delete from queue", "error", err)
				return
			}
		}
	}
}
