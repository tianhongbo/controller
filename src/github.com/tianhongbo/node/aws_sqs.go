package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

const AWS_SQS_URL = "https://sqs.us-west-2.amazonaws.com/497100832806/MTaaS-emulator-queue"
const AWS_SQS_ARN = "arn:aws:sqs:us-west-2:497100832806:MTaaS-emulator-queue"

const RSP_EMULATOR_URI = "http://mtaas-worker.us-west-2.elasticbeanstalk.com/api/v1/emulator/"

const AWS_CURRENT_REGION = "us-west-2"

const AWS_PUBLIC_IP_URI = "http://instance-data/latest/meta-data/public-ipv4"
const AWS_LOCAL_IP_URI = "http://instance-data/latest/meta-data/local-ipv4"
const AWS_AVAILABILITY_ZONE_URI = "http://instance-data/latest/meta-data/placement/availability-zone"

const SYS_LOG_FILE = "/home/ubuntu/controller/log/sys.log"

type LocalEnv struct {
	PublicIp         string
	LocalIp          string
	AvailabilityZone string
}

var localEnv LocalEnv

type Spec_T struct {
	Ram       string `json:"ram"`
	Cpu       string `json:"cpu"`
	OsVersion string `json:"os_version"`
}

type ReqEmulatorMsg struct {
	Id       string `json:"_id"`
	Status   string `json:"status"`
	Region   string `json:"region"`
	Occupant string `json:"occupant"`
	Adb_uri  string `json:"adb_uri"`
	spec     struct {
		Ram       string `json:"ram"`
		Cpu       string `json:"cpu"`
		OsVersion string `json:"os_version"`
	} `json:"spec"`
}

type RspEmulatorMsg struct {
}

type EmulatorNew struct {
	Name string
	Port int
	Used bool
}

// need to update with the maximum emulators supported
var emulatorsNew = []EmulatorNew{
	//1
	EmulatorNew{
		Name: "emulator-5554",
		Port: 5555,
		Used: false,
	},
	//2
	EmulatorNew{
		Name: "emulator-5556",
		Port: 5557,
		Used: false,
	},
	//3
	EmulatorNew{
		Name: "emulator-5558",
		Port: 5559,
		Used: false,
	},
	/*
		//4
		EmulatorNew {
			Name: "emulator-5560",
			Port: 5561,
			Used: false,
		},
		//5
		EmulatorNew {
			Name: "emulator-5562",
			Port: 5563,
			Used: false,
		},
		//6
		EmulatorNew {
			Name: "emulator-5564",
			Port: 5565,
			Used: false,
		},
		//7
		EmulatorNew {
			Name: "emulator-5566",
			Port: 5567,
			Used: false,
		},
		//8
		EmulatorNew {
			Name: "emulator-5568",
			Port: 5569,
			Used: false,
		},
		//9
		EmulatorNew {
			Name: "emulator-5570",
			Port: 5571,
			Used: false,
		},
		//10
		EmulatorNew {
			Name: "emulator-5572",
			Port: 5573,
			Used: false,
		},
	*/
}

func init() {
	//initialize log at the first place
	initLog()

	log.Println("system restart...")

	//initialize local env
	localEnv.PublicIp = getVmMetaData(AWS_PUBLIC_IP_URI)
	localEnv.LocalIp = getVmMetaData(AWS_LOCAL_IP_URI)
	localEnv.AvailabilityZone = getVmMetaData(AWS_AVAILABILITY_ZONE_URI)

	log.Println("local env initialization is done.")
	log.Println("pulic ip: ", localEnv.PublicIp)
	log.Println("local ip: ", localEnv.LocalIp)
	log.Println("availability zone: ", localEnv.AvailabilityZone)
}

func initLog() {
	file, err := os.OpenFile(SYS_LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("fail to open log file with error:  ", err)
	}

	multi := io.MultiWriter(file, os.Stdout)
	log.SetOutput(multi)
	log.Println("log initialization is done.")
}

func getVmMetaData(uri string) string {
	//get public ip
	resp, err := http.Get(uri)
	if err != nil {
		// handle error
		log.Println("fail to get public IP.")
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("fail to decode public ip.")
		return ""
	}
	return string(body)
}

func getOneEmulator() (EmulatorNew, error) {
	for i, v := range emulatorsNew {
		if v.Used != true {
			emulatorsNew[i].Used = true
			return v, nil
		}
	}
	e := EmulatorNew{}
	return e, errors.New("no available emulators.")
}

func hasFreeEmulator() bool {
	for _, v := range emulatorsNew {
		if v.Used {
			continue
		}
		return true
	}
	return false
}

func getOneMsg(svc *sqs.SQS) (ReqEmulatorMsg, string, error) {

	var req ReqEmulatorMsg

	params := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(AWS_SQS_URL), // Required
		AttributeNames: []*string{
			aws.String(AWS_SQS_ARN), // Required
			// More values...
		},
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(1),
	}

	resp, err := svc.ReceiveMessage(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		log.Println(err.Error())
		return ReqEmulatorMsg{}, "", errors.New("fail to unmarshal Json from create emulator request.")
	}
	if len(resp.Messages) < 1 {
		log.Println("no message")
		return ReqEmulatorMsg{}, "", errors.New("fail to unmarshal Json from create emulator request.")
	}
	body := []byte(*(resp.Messages[0].Body))

	// decode text message to json
	if err := json.Unmarshal(body, &req); err != nil {
		log.Println("Error! Can't unmarshal Json from create emulator request.")
		return ReqEmulatorMsg{}, "", errors.New("fail to unmarshal Json from create emulator request.")
	}

	// Pretty-print the response data.
	log.Println("receive one message with body: ", *(resp.Messages[0].Body))

	log.Println("receive one message with object: ", req.Id)

	return req, *(resp.Messages[0].ReceiptHandle), nil

}

// handle one message
func doOneMsg(req ReqEmulatorMsg) error {
	emu, err := getOneEmulator()
	if err != nil {
		// do not delete the message from sqs queue
		log.Println("fail to allocate a emulator.")
		return errors.New("fail to allocate a emlator.")
	}
	log.Println(req.Id, emu.Name)

	// reuse req to assemble rsp
	req.Status = "occupied"
	req.Adb_uri = localEnv.PublicIp + ":" + strconv.Itoa(emu.Port)

	// code to json
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(req)

	// set url
	url := RSP_EMULATOR_URI + req.Id

	// send PUT request
	client := &http.Client{}
	request, err := http.NewRequest("PUT", url, b)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return errors.New("fail to send PUT request.")
	}
	log.Println("send emulator req to: ", url)
	log.Println("receive emulator rsp: ", response)
	return nil
}

// delete one message
func deleteOneMsg(svc *sqs.SQS, ReceiptHandle string) {

	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(AWS_SQS_URL),   // Required
		ReceiptHandle: aws.String(ReceiptHandle), // Required
	}
	_, err := svc.DeleteMessage(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		log.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	log.Println("delete one message.")

}

func main() {

	sess, err := session.NewSession()
	if err != nil {
		log.Println("failed to create session,", err)
		return
	}

	svc := sqs.New(sess, &aws.Config{Region: aws.String(AWS_CURRENT_REGION)})

	for hasFreeEmulator() {
		// receive one message
		req, handler, err := getOneMsg(svc)
		if err != nil {
			continue
		}
		log.Println(req.Id, handler)

		// process the message
		err = doOneMsg(req)
		if err != nil {
			// do not delete the message from sqs queue
			log.Println("fail to process the request.")
			continue
		}

		// delete the message
		deleteOneMsg(svc, handler)

		// delay for
	}

	log.Println("exits because all emulators are used now.")
}
