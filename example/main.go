package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	spinhttp "github.com/fermyon/spin-go-sdk/http"
	"github.com/fermyon/spin-go-sdk/variables"

	s3 "github.com/fermyon/spin-s3-go"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		accessKeyID, err := variables.Get("aws_access_key_id")
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

		// aws config
		cfg := s3.Config{
			AccessKey:    accessKeyID,
			SecretKey:    secretAccessKey,
			SessionToken: sessionToken,
			Region:       "us-east-1",
			Endpoint:     "https://s3.us-east-1.amazonaws.com",
		}
		s3Client, err := s3.New(cfg)
		if err != nil {
			fmt.Println(err)
			return
		}

		ctx := context.Background()

		resp, err := s3Client.ListBuckets(ctx)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("-- List buckets --")
		for _, bucket := range resp.Buckets {
			fmt.Println(bucket)
		}

		const bucketName = "spin-s3-examples"
		const fileName = "hello.txt"

		fmt.Println("-- Creating object --")
		if err := s3Client.PutObject(ctx, bucketName, fileName, strings.NewReader("Hello S3!")); err != nil {
			fmt.Println(err)
			return
		}

		// Currently broken... =[
		fmt.Println("-- Getting object --")
		contents, err := s3Client.GetObject(ctx, bucketName, fileName)
		if err != nil {
			fmt.Println(err)
			return
		}

		b, err := io.ReadAll(contents)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))

		fmt.Println("Success")
	})
}

func main() {}
