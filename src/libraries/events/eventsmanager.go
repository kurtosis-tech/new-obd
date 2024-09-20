package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"os"
	"time"
)

type EventsManager struct {
	snsTopicARN string
	queueUrl    string
}

func newEventsManager(snsTopicARN string, queueUrl string) *EventsManager {
	return &EventsManager{snsTopicARN: snsTopicARN, queueUrl: queueUrl}
}

func CreateEventsManager() (*EventsManager, error) {

	//awsKey := os.Getenv("AWS_ACCESS_KEY_ID")
	//awsSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION")

	if awsRegion == "" {
		//return nil, errors.New("imposible to init the events manager component because the AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY or AWS_REGION environment variable not set")
		return nil, errors.New("imposible to init the events manager component because AWS_REGION environment variable not set")
	}

	snsTopicARN := os.Getenv("SNS_TOPIC_ARN")
	queueUrl := os.Getenv("QUEUE_URL")

	if snsTopicARN == "" || queueUrl == "" {
		return nil, errors.New("imposible to init the events manager component because the SNS_TOPIC_ARN or QUEUE_URL environment variable not set")
	}

	manager := newEventsManager(snsTopicARN, queueUrl)

	return manager, nil
}

func (manager *EventsManager) PublishMessage(message string) error {
	sess := session.Must(session.NewSession())
	svc := sns.New(sess)

	_, err := svc.Publish(&sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(manager.snsTopicARN),
	})
	return err
}

func (manager *EventsManager) ReceiveMessages(msgChan chan<- string, errChan chan<- error) {
	sess := session.Must(session.NewSession())
	svc := sqs.New(sess)

	for {
		result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(manager.queueUrl),
			MaxNumberOfMessages: aws.Int64(1),
			WaitTimeSeconds:     aws.Int64(10), // Long polling
		})
		if err != nil {
			errChan <- errors.New(fmt.Sprintf("error receiving message: %s", err))
		}

		for _, message := range result.Messages {

			var notification SNSNotification
			err = json.Unmarshal([]byte(*message.Body), &notification)
			if err != nil {
				errChan <- errors.New(fmt.Sprintf("Error unmarshalling JSON: %v\n", err))
			}

			msgChan <- notification.Message
			// Delete the message after processing
			_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      aws.String(manager.queueUrl),
				ReceiptHandle: message.ReceiptHandle,
			})
			if err != nil {
				errChan <- errors.New(fmt.Sprintf("error deleting message: %s", err))
			}
		}

		time.Sleep(1 * time.Second) // Prevent tight loop
	}

}
