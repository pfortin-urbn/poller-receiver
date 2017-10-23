package pollers

import (
	"os"
	"sync"
	"time"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
)

type AwsPoller struct {
	Wg            *sync.WaitGroup
	Done          chan bool
	Cs            chan string
	PollingPeriod int
	QueueName     string
	TimeoutChan   chan bool
}

func (poller AwsPoller) timeout() {
	time.Sleep(time.Duration(poller.PollingPeriod) * time.Millisecond)
	poller.TimeoutChan <- true
}

func (poller AwsPoller) GetMessages() {
	defer poller.Wg.Done()
	log.Println("Starting AWS Poller...")

	svc := sqs.New(session.New(), &aws.Config{Region: aws.String("us-east-1")})
	url := fmt.Sprintf("https://sqs.us-east-1.amazonaws.com/478989820108/%s", poller.QueueName)

	var waitTimeSecs int64 = 10
	params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(url), // Required
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds:     &waitTimeSecs,
	}

	//Set up the polling period timeout goroutine
	poller.TimeoutChan = make(chan bool, 1)
	go poller.timeout()

	for {
		select {
		case <-poller.Done:
			log.Debug("Aws Poller Done received.")
			return
		case <-poller.TimeoutChan:
			resp, err := svc.ReceiveMessage(params)
			if err != nil {
				log.Println(url)
				log.Printf("-->> " + err.Error())
				os.Exit(0)
			}

			if len(resp.Messages) > 0 {
				for _, msg := range resp.Messages {
					//log.Debug("Q-Url:", url, ", ReceiptHandle:", msg.ReceiptHandle)
					message := *msg.Body
					delParams := &sqs.DeleteMessageInput{
						QueueUrl:      aws.String(url),                // Required
						ReceiptHandle: aws.String(*msg.ReceiptHandle), // Required
					}
					poller.Cs <- message
					svc.DeleteMessage(delParams)
				}
			}
			go poller.timeout()
		}
	}
}
