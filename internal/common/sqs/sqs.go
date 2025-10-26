package sqs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	awsSQSTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"

	"github.com/lucaslobo/aggregator/internal/common/logs"
	"github.com/lucaslobo/aggregator/internal/inbound"
)

// client is a wrapper for the SQS client that makes working with SQS simpler
type client struct {
	logger              logs.Logger
	sqsClient           *sqs.Client
	sqsURL              string
	maxNumberOfMessages int
	waitTimeSeconds     int
}

// NewClient Creates a new SQS client wrapper
func NewClient(cfg ConfigSQS) inbound.Queue {
	return client{
		logger:              cfg.Logger,
		sqsClient:           cfg.SqsClient,
		sqsURL:              cfg.SqsURL,
		maxNumberOfMessages: cfg.MaxNumberOfMessages,
		waitTimeSeconds:     cfg.WaitTimeSeconds,
	}
}

// GetMessages fetches the messages from the queue
func (s client) GetMessages(ctx context.Context) (*sqs.ReceiveMessageOutput, error) {
	input := &sqs.ReceiveMessageInput{
		MaxNumberOfMessages: int32(s.maxNumberOfMessages),
		WaitTimeSeconds:     int32(s.waitTimeSeconds),
		QueueUrl:            aws.String(s.sqsURL),
	}

	msgOutput, err := s.sqsClient.ReceiveMessage(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("could not receive message from sqs: %w", err)
	}

	return msgOutput, nil
}

// SendMessage sends a message to the queue
func (s client) SendMessage(ctx context.Context, message awsSQSTypes.Message) error {
	_, err := s.sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: message.Body,
		QueueUrl:    aws.String(s.sqsURL),
	})
	if err != nil {
		return fmt.Errorf("could not send message to sqs: %w", err)
	}

	return nil
}

// Delete deletes a message from the queue
func (s client) Delete(ctx context.Context, message awsSQSTypes.Message) error {
	input := sqs.DeleteMessageInput{
		ReceiptHandle: message.ReceiptHandle,
		QueueUrl:      aws.String(s.sqsURL),
	}

	_, err := s.sqsClient.DeleteMessage(ctx, &input)
	if err != nil {
		return fmt.Errorf("could not delete message from sqs: %w", err)
	}
	return nil
}

// ChangeMessageVisibility changes the visibility timeout of a message in the queue
func (s client) ChangeMessageVisibility(ctx context.Context, receiptHandle *string, timeout int64) error {
	input := sqs.ChangeMessageVisibilityInput{
		QueueUrl:          aws.String(s.sqsURL),
		ReceiptHandle:     receiptHandle,
		VisibilityTimeout: int32(timeout),
	}

	_, err := s.sqsClient.ChangeMessageVisibility(ctx, &input)
	if err != nil {
		return fmt.Errorf("could not change message visibility in sqs: %w", err)
	}

	return nil
}
