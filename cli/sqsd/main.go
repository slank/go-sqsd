package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"
	sqsd "github.com/slank/go-sqsd"
)

type Config struct {
	QueueURL           string
	HTTPDestURL        string
	MessagesPerRequest int64
	PollWaitSeconds    int64
	SleepDuration      time.Duration
	ContentType        string
}

func main() {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	var conf Config

	fs.Int64Var(&conf.MessagesPerRequest, "sqs.messages_per_request", 10, "Maximum number of messages to receive per request")
	fs.Int64Var(&conf.PollWaitSeconds, "sqs.poll_wait_seconds", 20, "Long poll time in seconds")
	fs.DurationVar(&conf.SleepDuration, "sqs.sleep_duration", 10*time.Second, "After an empty receive, wait this long before next poll")
	fs.StringVar(&conf.ContentType, "http.content_type", "application/json", "Value of the Content-Type HTTP header")

	flag.Usage = func() {
		fmt.Printf("usage: %s [options] queue_url dest_url\n\n", os.Args[0])
		fs.PrintDefaults()
	}
	fs.Parse(os.Args[1:])

	args := fs.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}
	conf.QueueURL = args[0]
	conf.HTTPDestURL = args[1]

	msgs := make(chan *sqs.Message)
	del := make(chan *sqs.Message)

	handler := sqsd.NewSQSHandler(conf.QueueURL)
	handler.QueueURL = conf.QueueURL
	handler.MessagesPerRequest = conf.MessagesPerRequest
	handler.PollWaitSeconds = conf.PollWaitSeconds
	handler.SleepDuration = conf.SleepDuration

	go handler.Poller(msgs)
	go handler.Deleter(del)

	pusher := sqsd.NewHTTPPusher(conf.HTTPDestURL)
	pusher.ContentType = conf.ContentType
	pusher.Start(msgs, del)
}
