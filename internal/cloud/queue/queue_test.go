// +build local

package queue_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"

	"gitlab.unanet.io/devops/eve/internal/cloud/queue"
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
		Body:    "blah",
	}
	err := q.Message(&m)
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
	err := q.Delete("AQEBt6pj0R7OCNTC4BeCXzyWtiWvcB7sLv1HG8rohb5w2Qbuw22iS1sKiUcMveDgXbilP5SX0AmlaJPaItDzOC6Sp5GE2ANhuE83dMJ5trg2Numzuab9iwthAKbhYyF5YrJS2k1O3jO0GphIGjhAyIarGiHGxrR58+xaHs5EmacuUgG9i52FcSLvSePtFgNtJqJwuyVc/ikWPr5mBKkKqtjx0GfGHy7csrOUHk/JnmP7VZKXsB0mv0KAmEhg6H/nAxy/Y+rmVR2362j8+8FoNaGQsS+7Q61wd3Jh6Thk38Iy788=")
	require.NoError(t, err)
}
