package main

import (
	"encoding/json"
	"log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const AWS_SQS_URL = "https://sqs.us-west-2.amazonaws.com/497100832806/MTaaS-emulator-queue"
const AWS_SQS_ARN = "arn:aws:sqs:us-west-2:497100832806:MTaaS-emulator-queue"
const AWS_CURRENT_REGION = "us-west-2"

type ReqEmulatorMsg struct {
	Id		string	`json:"_id"`
	status		string  `json:"status"`
	region		string  `json:"region"`
	occupant	string  `json:"occupant"`
	adb_uri		string  `json:"adb_uri"`
	spec		string  `json:"spec"`
}

type RspEmulatorMsg struct {

}

func init() {
	//createEmulators()
}

func main() {
	var req ReqEmulatorMsg

	sess, err := session.NewSession()
	if err != nil {
		log.Println("failed to create session,", err)
		return
	}


	svc := sqs.New(sess, &aws.Config{Region: aws.String(AWS_CURRENT_REGION)})

	params := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(AWS_SQS_URL), // Required
		AttributeNames: []*string{
			aws.String(AWS_SQS_ARN), // Required
			// More values...
		},
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:       aws.Int64(1),
		WaitTimeSeconds:         aws.Int64(1),
	}
	resp, err := svc.ReceiveMessage(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		log.Println(err.Error())
		return
	}
	if len(resp.Messages) < 1 {
		log.Println("no message")
		return
	}
	body := []byte(*(resp.Messages[0].Body))

	// decode text message to json
        if err := json.Unmarshal(body, &req); err != nil {
                log.Println("Error! Can't unmarshal Json from create emulator request.")
                return
        }

	// Pretty-print the response data.
	log.Println("receive one message with body: ", *(resp.Messages[0].Body))

	log.Println("receive one message with object: ", req.Id)

	// handle the message

	// delete the message
	delete_params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(AWS_SQS_URL), // Required
		ReceiptHandle: aws.String(*(resp.Messages[0].ReceiptHandle)), // Required
	}
	_, err = svc.DeleteMessage(delete_params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		log.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	log.Println("delete one message.")

}
