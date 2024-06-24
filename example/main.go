package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	spinhttp "github.com/fermyon/spin-go-sdk/http"
	"github.com/fermyon/spin-go-sdk/variables"

	s3 "github.com/fermyon/spin-s3-go"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		endpoint, err := variables.Get("s3_endpoint")
		if err != nil {
			fmt.Println(err)
			return
		}

		cfg := s3.Config{
			// Setting config.Endpoint allows us to provide an endpoint other than
			// AWS. It is not required when using AWS S3.
			Endpoint: endpoint,

			// The following fields are for your AWS credentials.
			// AccessKey:    "",
			// SecretKey:    "",
			// SessionToken: "",
			// Region:       "",
		}

		s3Client, err := s3.New(cfg)
		if err != nil {
			fmt.Println(err)
			return
		}

		ctx := context.Background()

		const bucketName = "spin-s3-examples"

		fmt.Printf("-- Create bucket %q\n", bucketName)
		if err := s3Client.CreateBucket(ctx, bucketName); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("-- List all buckets")
		resp, err := s3Client.ListBuckets(ctx)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, bucket := range resp.Buckets {
			fmt.Println(bucket)
		}

		const objectName = "hello.txt"

		fmt.Printf("-- Creating object %q\n", objectName)
		if err := s3Client.PutObject(ctx, bucketName, objectName, []byte("Hello S3!")); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("-- Getting object")
		contents, err := s3Client.GetObject(ctx, bucketName, objectName)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Object contents:")
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
