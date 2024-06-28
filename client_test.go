package aws

import (
	"testing"
)

func TestClientBuildEndpoint(t *testing.T) {
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
		endpointURL: "https://s3.us-east-1.amazonaws.com",
		path:        "myobject",
		want:        "https://s3.us-east-1.amazonaws.com/myobject",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &S3Client{CommonClient{
				endpointURL: tt.endpointURL,
			}}
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
