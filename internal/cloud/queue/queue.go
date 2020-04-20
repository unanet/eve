package queue

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var (
	qURL = "https://sqs.us-east-2.amazonaws.com/580107804399/devops-prod-eve.fifo"
)

func ListQueues() (interface{}, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2")},
	)

	// Create a SQS service client.
	svc := sqs.New(sess)

	// List the queues available in a given region.
	result, err := svc.ListQueues(nil)
	if err != nil {
		fmt.Println("Error", err)
		return nil, err
	}

	fmt.Println("Success")
	// As these are pointers, printing them out directly would not be useful.
	for i, urls := range result.QueueUrls {
		// Avoid dereferencing a nil pointer.
		if urls == nil {
			continue
		}
		fmt.Printf("%d: %s\n", i, *urls)
	}

	return nil, nil
}

func CreateMessage(groupID, message string, id string) (*sqs.SendMessageOutput, error) {
	s, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}
	svc := sqs.New(s)

	result, err := svc.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"ID": {
				DataType:    aws.String("String"),
				StringValue: aws.String(id),
			},
		},
		MessageBody:            aws.String(message),
		MessageGroupId:         aws.String(groupID),
		MessageDeduplicationId: aws.String(id),
		QueueUrl:               &qURL,
	})

	if err != nil {
		fmt.Println("Error", err)
		return nil, err
	}

	fmt.Println("Success", *result.MessageId)
	return result, nil
}

func ReceiveMessage() (*sqs.ReceiveMessageOutput, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
			aws.String(sqs.MessageSystemAttributeNameMessageGroupId),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &qURL,
		MaxNumberOfMessages: aws.Int64(1),
		//VisibilityTimeout:   aws.Int64(20), // 20 seconds
		WaitTimeSeconds: aws.Int64(0),
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func DeleteMessage(receiptHandle string) (*sqs.DeleteMessageOutput, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)
	resultDelete, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &qURL,
		ReceiptHandle: aws.String(receiptHandle),
	})

	if err != nil {
		return nil, err
	}

	return resultDelete, nil
}
