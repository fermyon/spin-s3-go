package aws

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	spinhttp "github.com/fermyon/spin-go-sdk/http"
)

// TODO: @asteurer can the Config and Client structs be moved to types.go?
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

// TODO: Pretty sure this is overcomplicating it
// Client provides an interface for interacting with the S3 API.
type S3Client struct {
	// config      Config
	// endpointURL string
	CommonClient
}

// New creates a new Client.
func NewS3(config Config) (*S3Client, error) {
	u, err := url.Parse(config.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse endpoint: %w", err)
	}
	client := &S3Client{CommonClient{
		config:      config,
		endpointURL: u.String(),
	}}

	return client, nil
}

// buildEndpoint returns a endpoint
func (c *S3Client) buildEndpoint(bucketName, path string) (string, error) {
	u, err := url.Parse(c.endpointURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse endpoint: %w", err)
	}
	if bucketName != "" {
		u.Host = bucketName + "." + u.Host
	}
	return u.JoinPath(path).String(), nil

}

func (c *S3Client) newRequest(ctx context.Context, method, bucketName, path string, body []byte) (*http.Request, error) {
	endpointURL, err := c.buildEndpoint(bucketName, path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, endpointURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var awsDate AwsDate
	awsDate.time = time.Now()

	// Set the AWS authentication headers
	payloadHash := GetPayloadHash(body)
	// Removing the 'protocol://' and '/uri-path' from the host header
	req.Header.Set("host", strings.Split(endpointURL, "/")[2])
	req.Header.Set("content-length", fmt.Sprintf("%d", len(body)))
	req.Header.Set("x-amz-content-sha256", payloadHash)
	req.Header.Set("x-amz-date", awsDate.getTime())
	req.Header.Set("x-amz-security-token", c.config.SessionToken)
	// Optional
	req.Header.Set("user-agent", "spin-s3")
	req.Header.Set("authorization", GetAuthorizationHeader(&c.CommonClient, req, &awsDate, payloadHash))

	return req, nil
}

// do sends the request and handles any error response.
func (c *S3Client) do(req *http.Request) (*http.Response, error) {
	resp, err := spinhttp.Send(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Only checking for a status of 200 feels too specific.
	if resp.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		if err := xml.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		return nil, errorResponse
	}
	return resp, nil
}

func (c *S3Client) CreateBucket(ctx context.Context, name string) error {
	req, err := c.newRequest(ctx, http.MethodPut, "", name, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	resp.Body.Close()
	return err
}

// ListBuckets returns a list of buckets.
func (c *S3Client) ListBuckets(ctx context.Context) (*ListBucketsResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "", "", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var results ListBucketsResponse
	if err := xml.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &results, nil
}

// ListObjects returns a list of objects within a specified bucket.
func (c *S3Client) ListObjects(ctx context.Context, bucketName string) (*ListObjectsResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, bucketName, "", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("RAW RESULTS:")
	fmt.Println(resp)

	var results ListObjectsResponse
	if err := xml.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &results, nil
}

// PutObject uploads an object to the specified bucket.
func (c *S3Client) PutObject(ctx context.Context, bucketName, objectName string, data []byte) error {
	req, err := c.newRequest(ctx, http.MethodPut, bucketName, objectName, data)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// GetObject fetches an object from the specified bucket.
// TODO: Create a struct to contain meta? etag,last modified, etc
func (c *S3Client) GetObject(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	req, err := c.newRequest(ctx, http.MethodGet, bucketName, objectName, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	// It's the callers responsibility to close the reader.
	// defer resp.Body.Close()

	return resp.Body, nil
}

// DeleteObject deletes an object from the specified bucket.
func (c *S3Client) DeleteObject(ctx context.Context, bucketName, objectName string) error {
	req, err := c.newRequest(ctx, http.MethodDelete, bucketName, objectName, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
