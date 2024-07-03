package s3

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	aws "github.com/fermyon/spin-aws-go"
)

func TestS3ClientBuildEndpoint(t *testing.T) {
	tests := []struct {
		name        string
		endpointURL string
		bucketName  string
		path        string
		want        string
	}{{

		name:       "with a bucket provided",
		bucketName: "kickit",
		want:       "https://kickit.s3.us-east-1.amazonaws.com",
	}, {
		name: "without a bucket provided",
		want: "https://s3.us-east-1.amazonaws.com",
	}, {
		name: "with a path provided",
		path: "myobject",
		want: "https://s3.us-east-1.amazonaws.com/myobject",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := NewS3(aws.Config{Service: "s3", Region: "us-east-1"})
			got, err := c.buildEndpoint(tt.bucketName, tt.path)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestS3NewRequest(t *testing.T) {
	buildReq := func(req *http.Request, config aws.Config, awsDate aws.AwsDate, url string, payload []byte) {
		payloadHash := aws.GetPayloadHash(payload)
		headers := map[string]string{
			"host":                 strings.Split(url, "/")[2],
			"content-length":       fmt.Sprintf("%d", len(payload)),
			"x-amz-content-sha256": payloadHash,
			"x-amz-date":           awsDate.GetTime(),
			"x-amz-security-token": config.SessionToken,
			"user-agent":           "spin-s3",
			"authorization":        aws.GetAuthorizationHeader(&config, req, &awsDate, payloadHash),
		}
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	s3Client, _ := NewS3(aws.Config{
		AccessKeyId:     "accesskeyid",
		SecretAccessKey: "secretaccesskey",
		SessionToken:    "sessiontoken",
		Region:          "us-east-1",
		Service:         "s3",
	})

	awsDate := aws.AwsDate{Time: time.Date(2024, 07, 01, 0, 0, 0, 0, time.UTC)}

	// Test Params: Create Bucket
	createBucketName := "example-bucket"
	createBucketPayload := []byte("")
	createBucketUrl, _ := s3Client.buildEndpoint("", createBucketName)
	createBucketReq, _ := http.NewRequestWithContext(context.Background(), "PUT", createBucketUrl, bytes.NewReader(createBucketPayload))
	buildReq(createBucketReq, s3Client.config, awsDate, createBucketUrl, createBucketPayload)

	//newRequest(ctx context.Context, method string, bucketName string, path string, body []byte)
	tests := []struct {
		name       string
		client     *S3Client
		method     string
		bucketName string
		path       string
		body       []byte
		want       *http.Request
	}{{
		name:       "Create Bucket",
		client:     s3Client,
		method:     "PUT",
		bucketName: "",
		path:       "test-bucket",
		body:       []byte(""),
		want:       createBucketReq,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s3Client.newRequest(context.Background(), tt.method, tt.bucketName, tt.path, tt.body)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}
