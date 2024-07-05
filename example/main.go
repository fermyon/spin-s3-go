package main

import (
	"context"
	"fmt"
	"io"
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

		region, err := variables.Get("aws_region")
		if err != nil {
			fmt.Println(err)
			return
		}

		endpoint, err := variables.Get("aws_endpoint")
		if err != nil {
			fmt.Println(err)
			return
		}

		service, err := variables.Get("aws_service")
		if err != nil {
			fmt.Println(err)
			return
		}

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
			bucketName := "test-bucket"

			s3Client, err := s3.NewS3(cfg)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = s3Client.CreateBucket(ctx, bucketName)
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
			defer resp.Close()

			bytes, err := io.ReadAll(resp)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(string(bytes))
		} else if service == "sqs" {

			sqsClient, err := sqs.NewSQS(cfg)
			if err != nil {
				fmt.Println(err)
				return
			}

			createParams := sqs.CreateQueueParams{
				QueueName: "test-queue",
			}

			createResp, err := sqsClient.CreateQueue(ctx, createParams)
			if err != nil {
				fmt.Println(err)
				return
			}

			endpoint = createResp.QueueURL

			sendParams := sqs.SendMessageParams{
				QueueURL:    endpoint,
				MessageBody: "Hello, SQS!",
			}

			_, err = sqsClient.SendMessage(ctx, sendParams)
			if err != nil {
				fmt.Println(err)
				return
			}

			recParams := sqs.ReceiveMessageParams{
				QueueURL: endpoint,
			}

			recResp, err := sqsClient.ReceiveMessage(ctx, recParams)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(recResp.Messages[0].Body)
		}
	})
}

func main() {}
