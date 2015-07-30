package sqsd

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/service/sqs"
)

// An HTTPPusher posts sqs.Messages to an HTTP endpoint.
type HTTPPusher struct {
	URL         string
	ContentType string
}

// NewHTTPPusher creates a new HTTPPusher using default values.
func NewHTTPPusher(url string) *HTTPPusher {
	return &HTTPPusher{
		URL:         url,
		ContentType: "application/json",
	}
}

// Start processes sqs.Messages from msgs and forwards them to del if they
// are sent successfully to the HTTP endpoint.
func (h *HTTPPusher) Start(msgs chan *sqs.Message, del chan *sqs.Message) {
	for msg := range msgs {
		if body, err := json.Marshal(msg); err != nil {
			log.Printf("Error marshaling message: %s", err)
		} else {
			resp, err := http.Post(h.URL, h.ContentType, bytes.NewBuffer(body))
			if err != nil {
				log.Printf("Error delivering message: %s", err)
				continue
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				log.Printf("Got status %d from the HTTP server", resp.StatusCode)
				continue
			}
			del <- msg
		}
	}
}
