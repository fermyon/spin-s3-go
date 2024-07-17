package s3

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestGetAuthorizationHeader(t *testing.T) {
	setRequestHeaders := func(req *http.Request, headers map[string]string) {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	s3PutPayload := []byte("Hello, S3!")
	s3PutReq, err := http.NewRequest("PUT", "https://example-bucket.s3.us-west-2.amazonaws.com/test", bytes.NewReader(s3PutPayload))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	s3PutPayloadHash := getPayloadHash(s3PutPayload)
	s3PutConfig := Config{
		AccessKey:    "accesskey",
		SecretKey:    "secretaccesskey",
		SessionToken: "sessiontoken",
		Region:       "us-east-1",
	}

	now := time.Date(2024, 07, 01, 0, 0, 0, 0, time.UTC)

	setRequestHeaders(s3PutReq,
		map[string]string{
			"host":                 "example-bucket.s3.us-west-2.amazonaws.com",
			"content-length":       fmt.Sprintf("%d", len(s3PutPayload)),
			"x-amz-content-sha256": s3PutPayloadHash,
			"x-amz-date":           now.Format("20060102T150405Z"),
			"x-amz-security-token": s3PutConfig.SessionToken,
			"user-agent":           "spin-s3",
		})

	tests := []struct {
		name        string
		config      *Config
		req         *http.Request
		date        *time.Time
		payloadHash string
		want        string
	}{{
		name:        "S3 PUT request",
		config:      &s3PutConfig,
		req:         s3PutReq,
		date:        &now,
		payloadHash: s3PutPayloadHash,
		want:        "AWS4-HMAC-SHA256 Credential=accesskey/20240701/us-east-1/s3/aws4_request, SignedHeaders=host;x-amz-content-sha256;x-amz-date;x-amz-security-token, Signature=fda1dfb8cf3b1af0a7020ede35e48f9c4a4df8c5a5da32118289af1ec49ab34a",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getAuthorizationHeader(tt.req, tt.payloadHash, tt.config.Region, tt.config.AccessKey, tt.config.SessionToken, tt.config.SecretKey, *tt.date)
			if got != tt.want {
				t.Errorf("got: %v\nwant: %v", got, tt.want)
			}
		})
	}
}
