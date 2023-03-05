package interactors

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	awsSQSTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"

	"github.com/lucaslobo/aggregator-cli/internal/common/logs"
	"github.com/lucaslobo/aggregator-cli/internal/core/application"
	"github.com/lucaslobo/aggregator-cli/internal/core/domain"
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
	app         application.Application
}

func NewQueueConsumer(logger logs.Logger, queueClient Queue, app application.Application) QueueConsumer {
	return QueueConsumer{
		logger:      logger,
		queueClient: queueClient,
		app:         app,
	}
}

// PollAndProcess polls the queue and processes the messages
func (q *QueueConsumer) PollAndProcess(ctx context.Context, windowSize int) {
	q.app.Init(windowSize)

	// indefinitely poll queue for messages
	for {
		messages, err := q.readQueueMessages(ctx)
		if errors.Is(err, errNoMessages) {
			q.logger.Info("no messages found")
			continue
		} else if err != nil {
			q.logger.Errorw("unexpected error when reading queue",
				"error", err)
			continue
		}
		q.processMessages(ctx, messages)
	}
}

var errNoMessages = errors.New("no sqs messages found")

func (q *QueueConsumer) readQueueMessages(ctx context.Context) ([]awsSQSTypes.Message, error) {
	res, err := q.queueClient.GetMessages(ctx)
	if err != nil {
		return nil, err
	}
	if res == nil || len(res.Messages) == 0 {
		return nil, errNoMessages
	}

	return res.Messages, nil
}

func (q *QueueConsumer) processMessages(ctx context.Context, messages []awsSQSTypes.Message) {
	q.logger.Infow("read messages from queue", "quantity", len(messages))
	for _, message := range messages {
		if message.Body != nil {
			var event domain.TranslationDelivered
			err := json.Unmarshal([]byte(*message.Body), &event)
			if err != nil {
				q.logger.Errorw("error unmarshalling message", "error", err)
				return
			}

			err = q.app.ProcessEvent(event)
			if err != nil {
				q.logger.Errorw("could not process message", "error", err)
				return
			}

			err = q.queueClient.Delete(ctx, message)
			if err != nil {
				q.logger.Errorw("could not delete from queue", "error", err)
				return
			}
		}
	}
}
