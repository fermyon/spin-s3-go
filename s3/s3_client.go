package s3

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	aws "github.com/fermyon/spin-aws-go"
	spinhttp "github.com/fermyon/spin-go-sdk/http"
)

// Client provides an interface for interacting with the S3 API.
type Client struct {
	config      aws.Config
	endpointURL string
}

// New creates a new Client.
func NewS3(config aws.Config) (*Client, error) {
	u, err := url.Parse(config.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse endpoint: %w", err)
	}

	client := &Client{
		config:      config,
		endpointURL: u.String(),
	}

	return client, nil
}

// buildEndpoint returns an endpoint
func (c *Client) buildEndpoint(bucketName, path string) (string, error) {
	u, err := url.Parse(c.endpointURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse endpoint: %w", err)
	}
	if bucketName != "" {
		u.Host = bucketName + "." + u.Host
	}
	return u.JoinPath(path).String(), nil
}

func (c *Client) newRequest(ctx context.Context, method, bucketName, path string, body []byte) (*http.Request, error) {
	endpointURL, err := c.buildEndpoint(bucketName, path)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(endpointURL)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, endpointURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var awsDate aws.AwsDate
	awsDate.Time = time.Now()

	// Set the AWS authentication headers
	payloadHash := getPayloadHash(body)
	req.Header.Set("host", u.Host)
	req.Header.Set("content-length", fmt.Sprintf("%d", len(body)))
	req.Header.Set("x-amz-content-sha256", payloadHash)
	req.Header.Set("x-amz-date", awsDate.GetTime())
	req.Header.Set("x-amz-security-token", c.config.SessionToken)
	// Optional
	req.Header.Set("user-agent", "spin-s3")
	req.Header.Set("authorization", aws.GetAuthorizationHeader(&c.config, req, &awsDate, payloadHash))

	return req, nil
}

func getPayloadHash(payload []byte) string {
	hash := sha256.New()
	hash.Write(payload)
	return hex.EncodeToString(hash.Sum(nil))
}

// do sends the request and handles any error response.
func (c *Client) do(req *http.Request) (*http.Response, error) {
	resp, err := spinhttp.Send(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Only checking for a status of 200 feels too specific.
	if resp.StatusCode != http.StatusOK {
		var errorResponse aws.ErrorResponse
		if err := xml.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		return nil, errorResponse
	}
	return resp, nil
}

func (c *Client) CreateBucket(ctx context.Context, name string) error {
	req, err := c.newRequest(ctx, http.MethodPut, "", name, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	resp.Body.Close()
	return err
}

// ListBuckets returns a list of buckets.
func (c *Client) ListBuckets(ctx context.Context) (*ListBucketsResponse, error) {
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
func (c *Client) ListObjects(ctx context.Context, bucketName string) (*ListObjectsResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, bucketName, "", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var results ListObjectsResponse
	if err := xml.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &results, nil
}

// PutObject uploads an object to the specified bucket.
func (c *Client) PutObject(ctx context.Context, bucketName, objectName string, data []byte) error {
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
func (c *Client) GetObject(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
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
func (c *Client) DeleteObject(ctx context.Context, bucketName, objectName string) error {
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
