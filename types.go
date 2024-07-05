package aws

import (
	"fmt"
	"time"
)

// Config contains the available options for configuring a Client.
type Config struct {
	// AWS Access key ID
	AccessKeyId string
	// AWS Secret Access key
	SecretAccessKey string
	// AWS Session Token
	SessionToken string
	// AWS region
	Region string
	// AWS Service
	Service string
	// Endpoint is an optional override URL to the s3 service.
	Endpoint string
}

type AwsDate struct {
	Time time.Time
}

func (d *AwsDate) GetDate() string {
	return d.Time.UTC().Format("20060102")
}

func (d *AwsDate) GetTime() string {
	return d.Time.UTC().Format("20060102T150405Z")
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// https://docs.aws.amazon.com/AmazonS3/latest/API/ErrorResponses.html#RESTErrorResponses
type ErrorResponse struct {
	Code      string `xml:"Code"`
	Message   string `xml:"Message"`
	Resource  string `xml:"Resource"`
	RequestID string `xml:"RequestId"`
}
