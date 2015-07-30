# go-sqsd

Go-sqsd is a clone of the sqsd found on Worker applications in AWS's Elastic Beanstalk PaaS. Given a queue URL and an HTTP endpoint URL, it polls the queue for messages and POSTs them to the endpoint. If the endpoint responds with a 200 response, the messasge is deleted.

## Build/Install

To build the cli:

```
go install github.com/slank/go-sqsd/cli/sqsd
```

## Usage

```
usage: sqsd [options] queue_url dest_url

  -http.content_type="application/json": Value of the Content-Type HTTP header
  -sqs.messages_per_request=10: Maximum number of messages to receive per request
  -sqs.poll_wait_seconds=20: Long poll time in seconds
  -sqs.sleep_duration=10s: After an empty receive, wait this long before next poll
```

AWS credentials are passed via the standard methods (environment variables, user credential files, instance profiles, etc). The AWS CLI docs have a [concise description](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#config-settings-and-precedence) of the various sources. `sqsd` offers no command-line options for passing AWS credentials.

You may need to set the AWS_REGION environment variable to the region where the SQS queue is defined.
