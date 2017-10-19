package pollers

import (
	"sync"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
)

type AwsPoller struct {
	Wg sync.WaitGroup
	Cs chan string
}

func (poller AwsPoller)GetMessages() {
	defer poller.Wg.Done()
	log.Println("Starting AWS Poller...")

	svc := sqs.New(session.New(), &aws.Config{Region: aws.String("us-east-1")})
	url := "https://sqs.us-east-1.amazonaws.com/478989820108/PAUL_TEST"

	var waitTimeSecs int64 = 10
	params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(url), // Required
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds:     &waitTimeSecs,
	}

	for {
		resp, err := svc.ReceiveMessage(params)
		if err != nil {
			log.Println(url)
			log.Printf("-->> " + err.Error())
			os.Exit(0)
		}

		if len(resp.Messages) > 0 {
			for _, msg := range resp.Messages {
				log.Debug("Q-Url:", url, ", ReceiptHandle:", msg.ReceiptHandle)
				message := *msg.Body
				delParams := &sqs.DeleteMessageInput{
					QueueUrl:      aws.String(url),                // Required
					ReceiptHandle: aws.String(*msg.ReceiptHandle), // Required
				}
				svc.DeleteMessage(delParams)
				poller.Cs <- message
			}
		}
	}

}