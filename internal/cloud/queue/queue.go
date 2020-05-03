package queue

import (
	"context"
	"strings"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

const (
	MessageAttributeReqID string = "eve_req_id"
	MessageAttributeState string = "eve_state"
)

var (
	queueID uint64
)

type Q struct {
	aws  *sqs.SQS
	c    Config
	log  *zap.Logger
	id   uint64
	sess *session.Session
}

// Config struct
type Config struct {
	MaxNumberOfMessage int64  `split_words:"true" default:"10"`
	QueueURL           string `split_words:"true" required:"true"`
	WaitTimeSecond     int64  `split_words:"true" default:"20"`
	VisibilityTimeout  int64  `split_words:"true" default:"3600"`
}

func NewQ(sess *session.Session, config Config) *Q {
	qID := atomic.AddUint64(&queueID, 1)
	awsQ := sqs.New(sess)
	return &Q{
		id:   qID,
		c:    config,
		aws:  awsQ,
		sess: sess,
		log:  log.Logger.With(zap.String("queue_url", config.QueueURL), zap.Uint64("internal_queue_id", qID)),
	}
}

type M struct {
	ID            uuid.UUID
	ReqID         string
	GroupID       string
	Body          string
	ReceiptHandle string
	MessageID     string
	State         string
}

func (q *Q) logWith(m *M) *zap.Logger {
	return q.log.With(
		zap.String("message_id", m.ID.String()),
		zap.String("req_id", m.ReqID))
}

func (q *Q) Message(ctx context.Context, m *M) error {
	if len(m.State) == 0 {
		m.State = "empty"
	}
	awsM := sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			MessageAttributeReqID: {
				DataType:    aws.String("String"),
				StringValue: aws.String(m.ReqID),
			},
			MessageAttributeState: {
				DataType:    aws.String("String"),
				StringValue: aws.String(m.State),
			},
		},
		MessageGroupId:         aws.String(m.GroupID),
		MessageDeduplicationId: aws.String(m.ID.String()),
		QueueUrl:               &q.c.QueueURL,
	}

	if len(m.Body) > 0 {
		awsM.MessageBody = aws.String(m.Body)
	} else {
		awsM.MessageBody = aws.String(m.ID.String())
	}

	now := time.Now()
	result, err := q.aws.SendMessageWithContext(ctx, &awsM)
	if err != nil {
		return errors.Wrap(err)
	}
	elapsed := time.Since(now)
	q.logWith(m).Info("AWS SQS message sent", zap.Float64("elapsed_ms", float64(elapsed.Nanoseconds())/1000000.0))
	m.MessageID = *result.MessageId
	return nil
}

func (q *Q) Receive(ctx context.Context) ([]*M, error) {
	awsM := sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(MessageAttributeReqID),
			aws.String(sqs.MessageSystemAttributeNameMessageGroupId),
			aws.String(sqs.MessageSystemAttributeNameMessageDeduplicationId),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            aws.String(q.c.QueueURL),
		MaxNumberOfMessages: aws.Int64(q.c.MaxNumberOfMessage),
		VisibilityTimeout:   aws.Int64(q.c.VisibilityTimeout),
		WaitTimeSeconds:     aws.Int64(q.c.WaitTimeSecond),
	}
	result, err := q.aws.ReceiveMessageWithContext(ctx, &awsM)
	if err != nil {
		if strings.HasPrefix(err.Error(), "RequestCanceled") {
			return nil, nil
		}
		return nil, errors.Wrap(err)
	}

	var returnMs []*M
	for _, x := range result.Messages {
		id := uuid.FromStringOrNil(*x.Attributes[sqs.MessageSystemAttributeNameMessageDeduplicationId])
		m := M{
			ID:            id,
			GroupID:       *x.Attributes[sqs.MessageSystemAttributeNameMessageGroupId],
			ReqID:         *x.MessageAttributes[MessageAttributeReqID].StringValue,
			State:         *x.MessageAttributes[MessageAttributeState].StringValue,
			Body:          *x.Body,
			ReceiptHandle: *x.ReceiptHandle,
			MessageID:     *x.MessageId,
		}
		returnMs = append(returnMs, &m)
		q.logWith(&m).Info("AWS SQS message received")
	}

	return returnMs, nil
}

func (q *Q) Delete(ctx context.Context, m *M) error {
	now := time.Now()
	_, err := q.aws.DeleteMessageWithContext(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(q.c.QueueURL),
		ReceiptHandle: aws.String(m.ReceiptHandle),
	})
	if err != nil {
		return errors.Wrap(err)
	}
	elapsed := time.Since(now)
	q.logWith(m).Info("AWS SQS message deleted", zap.Float64("elapsed_ms", float64(elapsed.Nanoseconds())/1000000.0))
	return nil
}
