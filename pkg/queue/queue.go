package queue

import (
	"context"
	goerrors "github.com/pkg/errors"
	"strings"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"

	"github.com/unanet/go/pkg/errors"
	"github.com/unanet/go/pkg/json"
	"github.com/unanet/go/pkg/log"
)

const (
	MessageAttributeReqID   string = "x_req_id"
	MessageAttributeCommand string = "eve_cmd"
	MessageAttributeID      string = "eve_id"
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
	GroupID       string
	Body          json.Object
	ReceiptHandle string
	MessageID     string
	Command       string
	DedupeID      string
}

type mContext struct {
	M
	ctx context.Context
}

func (q *Q) logWith(ctx context.Context) *zap.Logger {
	return q.log.With(zap.String("req_id", log.GetReqID(ctx)))
}

func (q *Q) Message(ctx context.Context, m *M) error {
	if len(m.Command) == 0 {
		m.Command = "empty"
	}
	awsM := sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			MessageAttributeReqID: {
				DataType:    aws.String("String"),
				StringValue: aws.String(log.GetReqID(ctx)),
			},
			MessageAttributeCommand: {
				DataType:    aws.String("String"),
				StringValue: aws.String(m.Command),
			},
			MessageAttributeID: {
				DataType:    aws.String("String"),
				StringValue: aws.String(m.ID.String()),
			},
		},
		MessageGroupId: aws.String(m.GroupID),
		QueueUrl:       &q.c.QueueURL,
	}

	if len(m.DedupeID) > 0 {
		awsM.MessageDeduplicationId = aws.String(m.DedupeID)
	}

	q.logWith(ctx).Info("preparing to send message to queue", zap.String("queue", *awsM.QueueUrl))

	if len(m.Body) > 0 {
		awsM.MessageBody = aws.String(m.Body.String())
	} else {
		awsM.MessageBody = aws.String(m.ID.String())
	}

	now := time.Now()
	result, err := q.aws.SendMessageWithContext(ctx, &awsM)
	if err != nil {
		return errors.Wrap(err)
	}
	if result == nil {
		return goerrors.New("nil sqs message response")
	}

	m.MessageID = *result.MessageId
	q.logWith(ctx).Info("AWS SQS message sent",
		zap.Float64("elapsed_ms", float64(time.Since(now).Nanoseconds())/1000000.0),
		zap.Any("id", m.ID),
		zap.String("message_id", m.MessageID),
	)
	return nil
}

func (q *Q) Receive(ctx context.Context) ([]*mContext, error) {
	awsM := sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameMessageGroupId),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            aws.String(q.c.QueueURL),
		MaxNumberOfMessages: aws.Int64(q.c.MaxNumberOfMessage),
		VisibilityTimeout:   aws.Int64(q.c.VisibilityTimeout),
		WaitTimeSeconds:     aws.Int64(q.c.WaitTimeSecond),
	}
	result, err := q.aws.ReceiveMessage(&awsM)
	if err != nil {
		if strings.HasPrefix(err.Error(), "RequestCanceled") {
			return nil, nil
		}
		return nil, errors.Wrap(err)
	}

	var returnMs []*mContext
	for _, x := range result.Messages {
		id := uuid.FromStringOrNil(*x.MessageAttributes[MessageAttributeID].StringValue)
		m := M{
			ID:            id,
			GroupID:       *x.Attributes[sqs.MessageSystemAttributeNameMessageGroupId],
			Command:       *x.MessageAttributes[MessageAttributeCommand].StringValue,
			Body:          json.Object(*x.Body),
			ReceiptHandle: *x.ReceiptHandle,
			MessageID:     *x.MessageId,
		}
		mctx := context.WithValue(ctx, log.RequestIDKey, *x.MessageAttributes[MessageAttributeReqID].StringValue)
		returnMs = append(returnMs, &mContext{
			M:   m,
			ctx: mctx,
		})
		q.logWith(mctx).Info("AWS SQS message received",
			zap.Any("id", m.ID),
			zap.String("message_id", m.MessageID),
		)
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
	q.logWith(ctx).Info("AWS SQS message deleted",
		zap.Float64("elapsed_ms", float64(time.Since(now).Nanoseconds())/1000000.0),
		zap.Any("id", m.ID),
		zap.String("message_id", m.MessageID),
	)
	return nil
}
