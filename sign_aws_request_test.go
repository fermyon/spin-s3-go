package aws

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
		AccessKeyId:     "accesskey",
		SecretAccessKey: "secretaccesskey",
		SessionToken:    "sessiontoken",
		Region:          "us-east-1",
		Service:         "s3",
	}
	s3PutDate := AwsDate{
		Time: time.Date(2024, 07, 01, 0, 0, 0, 0, time.UTC),
	}

	setRequestHeaders(s3PutReq,
		map[string]string{
			"host":                 "example-bucket.s3.us-west-2.amazonaws.com",
			"content-length":       fmt.Sprintf("%d", len(s3PutPayload)),
			"x-amz-content-sha256": s3PutPayloadHash,
			"x-amz-date":           s3PutDate.GetTime(),
			"x-amz-security-token": s3PutConfig.SessionToken,
			"user-agent":           "spin-s3",
		})

	tests := []struct {
		name        string
		config      *Config
		req         *http.Request
		date        *AwsDate
		payloadHash string
		want        string
	}{{
		name:        "S3 PUT request",
		config:      &s3PutConfig,
		req:         s3PutReq,
		date:        &s3PutDate,
		payloadHash: s3PutPayloadHash,
		want:        "AWS4-HMAC-SHA256 Credential=accesskey/20240701/us-east-1/s3/aws4_request, SignedHeaders=content-length;host;user-agent;x-amz-content-sha256;x-amz-date;x-amz-security-token, Signature=c46f099626fe27776d566b251f00470d0f6260d7945d4a8db1fef7ce8aa64e03",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetAuthorizationHeader(tt.config, tt.req, tt.date, tt.payloadHash)
			if got != tt.want {
				t.Errorf("got: %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPayloadHash(t *testing.T) {
	tests := []struct {
		name    string
		payload []byte
		want    string
	}{{
		name:    "Hello, world!",
		payload: []byte("Hello, world!"),
		want:    "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getPayloadHash(tt.payload)
			if got != tt.want {
				t.Errorf("got: %v, want %v", got, tt.want)
			}
		})
	}
}
