package s3

import (
	"testing"

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

		name:        "with a bucket provided",
		endpointURL: "https://s3.us-east-1.amazonaws.com",
		bucketName:  "kickit",
		want:        "https://kickit.s3.us-east-1.amazonaws.com",
	}, {
		name:        "without a bucket provided",
		endpointURL: "https://s3.us-east-1.amazonaws.com",
		want:        "https://s3.us-east-1.amazonaws.com",
	}, {
		name:        "with a path provided",
		path:        "myobject",
		endpointURL: "https://s3.us-east-1.amazonaws.com",
		want:        "https://s3.us-east-1.amazonaws.com/myobject",
	}, {
		name:        "bucket and path provided",
		path:        "myobject",
		endpointURL: "https://s3.us-east-1.amazonaws.com",
		bucketName:  "kickit",
		want:        "https://kickit.s3.us-east-1.amazonaws.com/myobject",
	}, {
		name:        "localstack",
		endpointURL: "http://s3.localhost.localstack.cloud:4566",
		path:        "test-object",
		bucketName:  "test-bucket",
		want:        "http://test-bucket.s3.localhost.localstack.cloud:4566/test-object",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewS3(aws.Config{Service: "s3", Region: "us-east-1", Endpoint: tt.endpointURL})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

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
