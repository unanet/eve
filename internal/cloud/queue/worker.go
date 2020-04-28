package queue

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

// HandlerFunc is used to define the Handler that is run on for each message
type HandlerFunc func(msg *sqs.Message) error

// HandleMessage wraps a function for handling sqs messages
func (f HandlerFunc) HandleMessage(msg *sqs.Message) error {
	return f(msg)
}

// Handler interface
type Handler interface {
	HandleMessage(msg *sqs.Message) error
}

// InvalidEventError struct
type InvalidEventError struct {
	event string
	msg   string
}

func (e InvalidEventError) Error() string {
	return fmt.Sprintf("[Invalid Event: %s] %s", e.event, e.msg)
}

// NewInvalidEventError creates InvalidEventError struct
func NewInvalidEventError(event, msg string) InvalidEventError {
	return InvalidEventError{event: event, msg: msg}
}

// Worker struct
type Worker struct {
	Config    *Config
	SqsClient sqsiface.SQSAPI
}

// Config struct
type Config struct {
	MaxNumberOfMessage int64  `split_words:"true" default:"10"`
	QueueURL           string `split_words:"true" required:"true"`
	WaitTimeSecond     int64  `split_words:"true" default:"20"`
}

func New(sess *session.Session, config *Config) (*Worker, error) {
	client := sqs.New(sess)

	return &Worker{
		Config:    config,
		SqsClient: client,
	}, nil
}

// Start starts the polling and will continue polling till the application is forcibly stopped
func (worker *Worker) Start(ctx context.Context, h Handler) {
	for {
		select {
		case <-ctx.Done():
			//log.Println("worker: Stopping polling because a context kill signal was sent")
			return
		default:
			//worker.Log.Debug("worker: Start Polling")

			params := &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(worker.Config.QueueURL), // Required
				MaxNumberOfMessages: aws.Int64(worker.Config.MaxNumberOfMessage),
				AttributeNames: []*string{
					aws.String("All"), // Required
				},
				WaitTimeSeconds: aws.Int64(worker.Config.WaitTimeSecond),
			}

			resp, err := worker.SqsClient.ReceiveMessage(params)
			if err != nil {
				log.Println(err)
				continue
			}
			if len(resp.Messages) > 0 {
				worker.run(h, resp.Messages)
			}
		}
	}
}

// poll launches goroutine per received message and wait for all message to be processed
func (worker *Worker) run(h Handler, messages []*sqs.Message) {
	numMessages := len(messages)
	//worker.Log.Info(fmt.Sprintf("worker: Received %d messages", numMessages))

	var wg sync.WaitGroup
	wg.Add(numMessages)
	for i := range messages {
		go func(m *sqs.Message) {
			// launch goroutine
			defer wg.Done()
			if err := worker.handleMessage(m, h); err != nil {
				//worker.Log.Error(err.Error())
			}
		}(messages[i])
	}

	wg.Wait()
}

func (worker *Worker) handleMessage(m *sqs.Message, h Handler) error {
	var err error
	err = h.HandleMessage(m)
	if _, ok := err.(InvalidEventError); ok {
		//worker.Log.Error(err.Error())
	} else if err != nil {
		return err
	}

	//worker.Log.Debug(fmt.Sprintf("worker: deleted message from queue: %s", aws.StringValue(m.ReceiptHandle)))

	return nil
}
