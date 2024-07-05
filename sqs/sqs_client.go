package sqs

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	aws "github.com/fermyon/spin-aws-go"
	spinhttp "github.com/fermyon/spin-go-sdk/http"
)

// Client provides an interface for interacting with the SQS API.
type Client struct {
	config      aws.Config
	endpointURL string
}

// New creates a new Client.
func NewSQS(config aws.Config) (*Client, error) {
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

func (c *Client) newSqsRequest(ctx context.Context, method, action string, body []byte) (*http.Request, error) {
	u, err := url.Parse(c.endpointURL)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var awsDate aws.AwsDate
	awsDate.Time = time.Now()

	// Set the AWS authentication headers
	payloadHash := getPayloadHash(body)
	req.Header.Set("host", u.Host)
	req.Header.Set("content-type", "application/x-amz-json-1.0")
	req.Header.Set("content-length", fmt.Sprintf("%d", len(body)))
	req.Header.Set("connection", "Keep-Alive")
	req.Header.Set("x-amz-target", action)
	req.Header.Set("x-amz-content-sha256", payloadHash)
	req.Header.Set("x-amz-date", awsDate.GetTime())
	req.Header.Set("x-amz-security-token", c.config.SessionToken)
	// Optional
	req.Header.Set("user-agent", "spin-sqs")
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

	// TODO: Only checking for a status of 200 feels too specific.
	if resp.StatusCode != http.StatusOK {
		var errorResponse aws.ErrorResponse
		if err := xml.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		return nil, errorResponse
	}
	return resp, nil
}

func (c *Client) CreateQueue(ctx context.Context, params CreateQueueParams) (*CreateQueueResponse, error) {
	messageBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode message parameters: %w", err)
	}

	req, err := c.newSqsRequest(ctx, http.MethodPost, "AmazonSQS.CreateQueue", messageBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response body: %w", err)
	}

	var sqsResponse CreateQueueResponse

	err = json.Unmarshal(bodyBytes, &sqsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode the response body %w", err)
	}

	return &sqsResponse, nil

}

func (c *Client) SendMessage(ctx context.Context, params SendMessageParams) (*SendMessageResponse, error) {
	messageBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode message parameters: %w", err)
	}

	req, err := c.newSqsRequest(ctx, http.MethodPost, "AmazonSQS.SendMessage", messageBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response body: %w", err)
	}

	var sqsResponse SendMessageResponse

	err = json.Unmarshal(bodyBytes, &sqsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode the response body %w", err)
	}

	return &sqsResponse, nil
}

func (c *Client) ReceiveMessage(ctx context.Context, params ReceiveMessageParams) (*ReceiveMessageResponse, error) {
	paramJsonBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode message parameters: %w", err)
	}

	req, err := c.newSqsRequest(ctx, http.MethodPost, "AmazonSQS.ReceiveMessage", paramJsonBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response body: %w", err)
	}

	var sqsResponse ReceiveMessageResponse

	err = json.Unmarshal(bodyBytes, &sqsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode the response body %w", err)
	}

	return &sqsResponse, nil
}

func (c *Client) DeleteMessage(ctx context.Context, params DeleteMessageParams) error {
	messageBytes, err := json.Marshal(params)
	if err != nil {
		return err
	}

	req, err := c.newSqsRequest(ctx, http.MethodPost, "AmazonSQS.DeleteMessage", messageBytes)
	if err != nil {
		return fmt.Errorf("failed to encode message parameters: %w", err)
	}

	_, err = c.do(req)

	return fmt.Errorf("failed to send SQS message: %w", err)
}
