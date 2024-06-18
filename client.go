package s3

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	spinhttp "github.com/fermyon/spin-go-sdk/http"
)

type Config struct {
	AccessKey    string
	SecretKey    string
	SessionToken string
	Region       string

	// Endpoint is the URL to the s3 service.
	Endpoint string
}

type Client struct {
	config     Config
	HTTPClient *http.Client

	// trace is only for dev. It needs to be removed to work within a spin app.
	trace bool
}

// New creates a new Client.
func New(config Config) (*Client, error) {
	client := &Client{
		config: config,

		// TODO: replace with spinhttp client.
		HTTPClient: spinhttp.NewClient(),
		// HTTPClient: http.DefaultClient,
	}

	return client, nil
}

func (c *Client) newRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "spin-s3")

	now := time.Now().UTC()

	// Set the AWS authentication headers
	req.Header.Set("Authorization", c.getAuthorizationHeader(req, now))
	req.Header.Set("x-amz-date", now.Format("20060102T150405Z"))
	req.Header.Set("x-amz-security-token", c.config.SessionToken)
	req.Header.Set("x-amz-content-sha256", getPayloadHash(""))

	return req, nil
}

// do is a temporary wrapper for making the request.
// This will need to be removed before running in a spin app because tinygo
// removed httputil. Just to be annoying.
func (c *Client) do(req *http.Request) (*http.Response, error) {
	// if c.trace {
	// 	b, err := httputil.DumpRequest(req, true)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	fmt.Println("-- REQUEST -------------------")
	// 	fmt.Println(string(b))
	// 	fmt.Println("-- END REQUEST ---------------")
	// }

	// There is a bug in tinygo where the Transport is not being called. As
	// a work around we need to call spinhttp.Send directly

	// resp, err := c.HTTPClient.Do(req)
	resp, err := spinhttp.Send(req)
	if err != nil {
		return nil, err
	}

	// if c.trace {
	// 	b, err := httputil.DumpResponse(resp, true)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	fmt.Println("-- RESPONSE ------------------")
	// 	fmt.Println(string(b))
	// 	fmt.Println("-- END RESPONSE --------------")
	// }

	return resp, nil
}

func (c *Client) CreateBucket(ctx context.Context, name string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.endpoint(name), nil)
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
	req, err := c.newRequest(ctx, http.MethodGet, c.endpoint(), strings.NewReader(""))
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
		return nil, err
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
		return nil, err
	}

	return &results, nil
}

// PutObject uploads an object to the specified bucket.
func (c *Client) PutObject(ctx context.Context, bucketName, objectName string, r io.Reader) error {
	req, err := c.newRequest(ctx, http.MethodPut, c.endpoint(bucketName, objectName), r)
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
	endpoint := c.endpoint(bucketName, objectName)
	req, err := c.newRequest(ctx, http.MethodGet, endpoint, nil)
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
	for len(endpoint) > 0 && endpoint[len(endpoint)-1] == '/' {
		endpoint = endpoint[0 : len(endpoint)-1]
	}
	return endpoint + "/" + path.Join(elem...)
}
