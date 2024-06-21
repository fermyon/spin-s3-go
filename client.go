package s3

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"time"

	spinhttp "github.com/fermyon/spin-go-sdk/http"
)

const (
	userAgent  = "spin-s3"
	timeFormat = "20060102T150405Z"
	dateFormat = "20060102"
)

// Config contains the available options for configuring a Client.
type Config struct {
	// S3 Access key ID
	AccessKey string
	// S3 Secret Access key
	SecretKey string
	// S3 Session Token
	SessionToken string
	// S3 region
	Region string

	// Endpoint is an optional override URL to the s3 service.
	Endpoint string
}

// validate checks for valid config options.
func (c *Config) validate() error {
	if c.Endpoint == "" && c.Region == "" {
		return errors.New("config Endpoint or Region must be set")
	}
	return nil
}

// Client provides an interface for interacting with the S3 API.
type Client struct {
	config     Config
	HTTPClient *http.Client
}

// New creates a new Client.
func New(config Config) (*Client, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}
	client := &Client{
		config:     config,
		HTTPClient: spinhttp.NewClient(),
		// HTTPClient: http.DefaultClient,
	}

	return client, nil
}

func (c *Client) newRequest(ctx context.Context, method, url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	now := time.Now().UTC()

	// Set the AWS authentication headers
	payloadHash := getPayloadHash(body)
	req.Header.Set("Authorization", getAuthorizationHeader(req, payloadHash, c.config.Region, c.config.AccessKey, c.config.SessionToken, c.config.SecretKey, now))
	req.Header.Set("x-amz-content-sha256", payloadHash)
	req.Header.Set("x-amz-date", now.Format(timeFormat))
	req.Header.Set("x-amz-security-token", c.config.SessionToken)

	req.Header.Set("User-Agent", userAgent)

	return req, nil
}

// do is a temporary wrapper for making the request.
func (c *Client) do(req *http.Request) (*http.Response, error) {
	// There is a bug in tinygo where the Transport is not being called. As
	// a work around we need to call spinhttp.Send directly
	// resp, err := c.HTTPClient.Do(req)
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

func (c *Client) CreateBucket(ctx context.Context, name string) error {
	req, err := c.newRequest(ctx, http.MethodPut, c.endpoint(name), nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}

	// TODO: Do I need to do anything with the response here?
	_ = resp

	return nil
}

// ListBuckets returns a list of buckets.
func (c *Client) ListBuckets(ctx context.Context) (*ListBucketsResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, c.endpoint(), nil)
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
	req, err := c.newRequest(ctx, http.MethodGet, c.endpoint(bucketName), nil)
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
	req, err := c.newRequest(ctx, http.MethodPut, c.endpoint(bucketName, objectName), data)
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

// GetObject fetches an object.
// TODO: Create a struct to contain meta? etag,last modified, etc
func (c *Client) GetObject(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	req, err := c.newRequest(ctx, http.MethodGet, c.endpoint(bucketName, objectName), nil)
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

func (c *Client) endpoint(elem ...string) string {
	endpoint := c.config.Endpoint
	// Strip trailing slashes.
	for endpoint != "" && endpoint[len(endpoint)-1] == '/' {
		endpoint = endpoint[0 : len(endpoint)-1]
	}
	// Build endpoint URL if no config.Endpoint is set.
	if endpoint == "" && c.config.Region != "" {
		endpoint = fmt.Sprintf("https://s3.%s.amazonaws.com", c.config.Region)
	}
	return endpoint + "/" + path.Join(elem...)
}
