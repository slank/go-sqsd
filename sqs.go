package sqsd

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// An SQSHandler can poll an SQS queue for messages and delete them after
// processing.
type SQSHandler struct {
	QueueURL           string
	MessagesPerRequest int64
	PollWaitSeconds    int64
	SleepDuration      time.Duration

	client *sqs.SQS
}

// NewSQSHandler creates an SQSHandler with default values.
func NewSQSHandler(queueURL string) *SQSHandler {
	sess := session.Must(session.NewSessionWithOptions(session.Options {
	        SharedConfigState: session.SharedConfigEnable,
        }))

	return &SQSHandler{
		QueueURL:           queueURL,
		MessagesPerRequest: 10,
		PollWaitSeconds:    20,
		SleepDuration:      10 * time.Second,
		client:             sqs.New(sess),
	}
}

// Poller begins polling the queue, pushing each received message to the
// channel provided.
func (h *SQSHandler) Poller(msgs chan *sqs.Message) {
	params := &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(h.QueueURL),
		AttributeNames:        []*string{aws.String("All")},
		MaxNumberOfMessages:   aws.Int64(h.MessagesPerRequest),
		MessageAttributeNames: []*string{aws.String("All")},
		WaitTimeSeconds:       aws.Int64(h.PollWaitSeconds),
	}
	for {
		received, err := h.client.ReceiveMessage(params)
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				log.Printf("Error reading queue: %s", awsErr)
				// TODO: update receive_error metric
				time.Sleep(h.SleepDuration)
			}
		} else {
			if len(received.Messages) == 0 {
				time.Sleep(h.SleepDuration)
			} else {
				for _, msg := range received.Messages {
					msgs <- msg
				}
			}
		}
	}
}

// Deleter deletes each sqs.Message sent to its channel
func (h *SQSHandler) Deleter(msgs chan *sqs.Message) {
	for msg := range msgs {
		_, err := h.client.DeleteMessage(
			&sqs.DeleteMessageInput{
				QueueUrl:      aws.String(h.QueueURL),
				ReceiptHandle: aws.String(*msg.ReceiptHandle),
			},
		)
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				log.Printf("Error deleting message: %s", awsErr)
				// TODO: update delete_error metric
			}
		}
	}
}
