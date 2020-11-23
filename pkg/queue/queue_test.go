// +build local

package queue_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"

	"gitlab.unanet.io/devops/eve/pkg/json"
	"gitlab.unanet.io/devops/eve/pkg/queue"
)

const (
	qURL = "https://sqs.us-east-2.amazonaws.com/580107804399/eve-api-prod.fifo"
)

func GetQueue(t *testing.T) *queue.Q {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2")},
	)
	require.NoError(t, err)
	return queue.NewQ(sess, queue.Config{
		MaxNumberOfMessage: 10,
		QueueURL:           qURL,
		WaitTimeSecond:     0,
		VisibilityTimeout:  3600,
	})
}

func TestQ_Message(t *testing.T) {
	q := GetQueue(t)
	m := queue.M{
		ID:      uuid.NewV4(),
		GroupID: "testing",
		Body:    json.Text("{\"blah:\",\"\"}"),
	}
	err := q.Message(context.TODO(), &m)
	require.NoError(t, err)
	fmt.Println(m)
}

func TestQ_Receive(t *testing.T) {
	q := GetQueue(t)
	ms, err := q.Receive()
	require.NoError(t, err)
	for _, x := range ms {
		fmt.Println(x)
	}
}

func TestQ_Delete(t *testing.T) {
	q := GetQueue(t)
	err := q.Delete(context.TODO(), &queue.M{
		ID:            uuid.UUID{},
		ReqID:         "",
		GroupID:       "",
		Body:          json.EmptyJSONText,
		ReceiptHandle: "",
		MessageID:     "",
		Command:       "",
	})
	require.NoError(t, err)
}
