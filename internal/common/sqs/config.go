package sqs

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/lucaslobo/aggregator-cli/internal/common/logs"
)

// ConfigSQS is used to provide configuration parameters to set up the SQS client wrapper
type ConfigSQS struct {
	Logger              logs.Logger
	SqsClient           sqs.Client
	SqsURL              string
	MaxNumberOfMessages int
	WaitTimeSeconds     int
}
