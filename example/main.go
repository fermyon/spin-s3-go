package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	s3 "github.com/fermyon/spin-s3-go"
	spinhttp "github.com/fermyon/spin/sdk/go/v2/http"
	"github.com/fermyon/spin/sdk/go/v2/variables"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		endpoint, err := variables.Get("s3_endpoint")
		if err != nil {
			fmt.Println("failed to get s3_endpoint variable:", err)
			return
		}

		// Create a Config with appropriate credentials.
		cfg := s3.Config{
			Endpoint: endpoint,

			// The following fields are for your AWS credentials.
			// AccessKey:    "",
			// SecretKey:    "",
			// SessionToken: "",
			// Region:       "",
		}

		// Create a New S3 client.
		s3Client, err := s3.New(cfg)
		if err != nil {
			fmt.Println("failed to create S3 client:", err)
			return
		}

		const bucketName = "spin-s3-examples"

		ctx := context.Background()
		fmt.Printf("-- Create bucket %q\n", bucketName)
		if err := s3Client.CreateBucket(ctx, bucketName); err != nil {
			fmt.Printf("failed to create bucket %q: %s\n", bucketName, err)
			return
		}

		fmt.Println("-- List all buckets")
		resp, err := s3Client.ListBuckets(ctx)
		if err != nil {
			fmt.Println("failed to list buckets:", err)
			return
		}
		for _, bucket := range resp.Buckets {
			fmt.Println(bucket)
		}

		const objectName = "hello.txt"

		fmt.Printf("-- Creating object %q\n", objectName)
		if err := s3Client.PutObject(ctx, bucketName, objectName, []byte("Hello S3!")); err != nil {
			fmt.Printf("failed to put object %q: %s\n", objectName, err)
			return
		}

		fmt.Printf("-- Getting object: %q\n", objectName)
		contents, err := s3Client.GetObject(ctx, bucketName, objectName)
		if err != nil {
			fmt.Printf("failed to get object %q: %s\n", objectName, err)
			return
		}

		fmt.Println("Object contents:")
		b, err := io.ReadAll(contents)
		if err != nil {
			fmt.Println("failed to read response:", err)
			return
		}
		fmt.Println(string(b))

		w.WriteHeader(http.StatusOK)
	})
}
