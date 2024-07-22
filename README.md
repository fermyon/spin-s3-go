# S3 Client for Spin Applications

[![Go Reference](https://pkg.go.dev/badge/github.com/fermyon/spin-s3-go.svg)](https://pkg.go.dev/github.com/fermyon/spin-s3-go)

This package provides an SDK for S3 compatible APIs for [Spin](https://developer.fermyon.com/spin/v2/index) applications.

See [example](./example) for a working example of how to use the SDK in your application.

## Usage

Add the client to you Spin application

```console
go get github.com/fermyon/spin-s3-go
```

### Example for creating an object

```go
package main

import (
	"context"
	"fmt"
	"net/http"

	spinhttp "github.com/fermyon/spin-go-sdk/http"

	s3 "github.com/fermyon/spin-s3-go"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		// Create a Config with appropriate credentials.
		cfg := s3.Config{
			AccessKey:    "your-access-key-id",
			SecretKey:    "your-secret-access-key",
			SessionToken: "your-session-token",
			Region:       "your-region",
		}

		// Create a New S3 client.
		s3Client, err := s3.New(cfg)
		if err != nil {
			fmt.Println("failed to create S3 client:", err)
			return
		}

		bucketName := "your-bucket-name"
		objectName := "greetings.txt"
		objectContents := []byte("Hello S3!")

		ctx := context.Background()
		if err := s3Client.PutObject(ctx, bucketName, objectName, objectContents); err != nil {
			fmt.Printf("failed to put object: %s\n", err)
		}
	})
}

func main() {}
```
