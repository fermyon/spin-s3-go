package aws

import (
	"net/http"
	"testing"
	"time"
)

func setRequestHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

// TODO: Come up with how a header should look first, then test
func TestGetAuthorizationHeader(t *testing.T) {
	s3GetReq, _ := http.NewRequest("GET", "https://test.s3.us-west-2.amazonaws.com/test", nil)
	s3GetPayloadHash := GetPayloadHash([]byte(""))
	s3GetConfig := Config{
		AccessKeyId:     "TEST",
		SecretAccessKey: "TEST",
		SessionToken:    "TEST",
		Region:          "us-east-1",
		Service:         "s3",
	}
	s3GetDate := AwsDate{
		Time: time.Date(2024, 07, 01, 0, 0, 0, 0, time.UTC),
	}

	setRequestHeaders(s3GetReq,
		map[string]string{
			"host":                 "test.s3.us-west-2.amazonaws.com",
			"content-length":       "0",
			"x-amz-content-sha256": s3GetPayloadHash,
			"x-amz-date":           s3GetDate.GetTime(),
			"x-amz-security-token": s3GetConfig.SessionToken,
			"user-agent":           "spin-s3",
		})

	tests := []struct {
		testName    string
		config      *Config
		req         *http.Request
		date        *AwsDate
		payloadHash string
		want        string
	}{{
		testName:    "S3 GET request",
		config:      &s3GetConfig,
		req:         s3GetReq,
		date:        &s3GetDate,
		payloadHash: s3GetPayloadHash,
		want:        "AWS4-HMAC-SHA256 Credential=TEST/20240701/us-east-1/s3/aws4_request, SignedHeaders=content-length;host;user-agent;x-amz-content-sha256;x-amz-date;x-amz-security-token, Signature=aec958a78ab4ce0c4cfe297558911531fe19ca9d47a59cf0d17a87eeb974fabd",
	}}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got := GetAuthorizationHeader(tt.config, tt.req, tt.date, tt.payloadHash)
			if got != tt.want {
				t.Errorf("got: %v, want %v", got, tt.want)
			}
		})
	}
}

// I used https://emn178.github.io/online-tools/sha256.html for the hmac generation testing
func TestGetPayloadHash(t *testing.T) {
	tests := []struct {
		testName string
		payload  []byte
		want     string
	}{{
		testName: "Hello, world!",
		payload:  []byte("Hello, world!"),
		want:     "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
	}}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got := GetPayloadHash(tt.payload)
			if got != tt.want {
				t.Errorf("got: %v, want %v", got, tt.want)
			}
		})
	}
}
