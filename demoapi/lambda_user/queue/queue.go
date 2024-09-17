package queue

import (
	"lambda-func/common"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type MessageQueue interface {
	SendMessage(message string) error
}

type SqsClient struct {
	sqsClient *sqs.SQS
}

func NewSqsClient() SqsClient {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := sqs.New(sess)

	return SqsClient{
		sqsClient: client,
	}
}

func (s SqsClient) SendMessage(message string) error {

	queueUrl, err := s.getQueueURL(common.QueueName)
	if err != nil {
		log.Printf("Failed to get queue URL: %v", err)
		return err
	}

	input := &sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    aws.String(queueUrl),
	}

	_, err = s.sqsClient.SendMessage(input)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}

	return nil

}

func (s SqsClient) getQueueURL(queueName string) (string, error) {
	input := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}

	result, err := s.sqsClient.GetQueueUrl(input)
	if err != nil {
		return "", err
	}

	return *result.QueueUrl, nil
}
