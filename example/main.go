package main

import (
	"context"
	"fmt"
	"net/http"

	spinhttp "github.com/fermyon/spin-go-sdk/http"
	"github.com/fermyon/spin-go-sdk/variables"

	aws "github.com/fermyon/spin-aws-go"
	s3 "github.com/fermyon/spin-aws-go/s3"
	sqs "github.com/fermyon/spin-aws-go/sqs"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		accessKeyId, err := variables.Get("aws_access_key_id")
		if err != nil {
			fmt.Println(err)
			return
		}

		secretAccessKey, err := variables.Get("aws_secret_access_key")
		if err != nil {
			fmt.Println(err)
			return
		}

		sessionToken, err := variables.Get("aws_session_token")
		if err != nil {
			fmt.Println(err)
			return
		}

		// Required header values
		endpoint := r.Header.Get("x-aws-endpoint")
		region := r.Header.Get("x-aws-region")
		service := r.Header.Get("x-aws-service")

		cfg := aws.Config{
			AccessKeyId:     accessKeyId,
			SecretAccessKey: secretAccessKey,
			SessionToken:    sessionToken,
			Endpoint:        endpoint,
			Region:          region,
			Service:         service,
		}

		ctx := context.Background()

		if service == "s3" {
			// S3 specific headers
			bucketName := r.Header.Get("x-s3-bucket")

			s3Client, err := s3.NewS3(cfg)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = s3Client.PutObject(ctx, bucketName, "hello.txt", []byte("Hello, S3!"))
			if err != nil {
				fmt.Println(err)
				return
			}

			resp, err := s3Client.GetObject(ctx, bucketName, "hello.txt")
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(resp)
		} else if service == "sqs" {
			queueUrl := r.Header.Get("x-sqs-queue-url")
			if queueUrl == "" {
				fmt.Println("If making an SQS request, you must include the 'x-sqs-queue-url' header.")
			}

			sqsClient, err := sqs.NewSQS(cfg)
			if err != nil {
				fmt.Println(err)
				return
			}

			sendParams := sqs.SqsSendMessageParams{
				QueueUrl:    queueUrl,
				MessageBody: "Hello, SQS!",
			}

			sendResp, err := sqsClient.SendMessage(ctx, sendParams)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(sendResp)

			recParams := sqs.SqsReceiveMessageParams{
				QueueUrl: queueUrl,
			}

			recResp, err := sqsClient.ReceiveMessage(ctx, recParams)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(recResp)
		}
	})
}

func main() {}
