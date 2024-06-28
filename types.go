package aws

import (
	"fmt"
	"time"
)

// https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListBuckets.html#API_ListBuckets_ResponseSyntax
type ListBucketsResponse struct {
	Buckets []BucketInfo `xml:"Buckets>Bucket"`
	Owner   Owner
}

// https://docs.aws.amazon.com/AmazonS3/latest/API/API_Bucket.html
type BucketInfo struct {
	Name         string
	CreationDate time.Time
}

// https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListObjects.html#API_ListObjects_ResponseSyntax
type ListObjectsResponse struct {
	CommonPrefixes []CommonPrefix
	Contents       []ObjectInfo
	Delimiter      string
	EncodingType   string
	IsTruncated    bool
	Marker         string
	MaxKeys        int
	Name           string
	NextMarker     string
	Prefix         string
}

// https://docs.aws.amazon.com/AmazonS3/latest/API/API_CommonPrefix.html
type CommonPrefix struct {
	Prefix string
}

// https://docs.aws.amazon.com/AmazonS3/latest/API/API_Object.html
type ObjectInfo struct {
	Key          string
	ETag         string
	Size         int
	LastModified time.Time
	StorageClass string
	Owner        Owner
}

// https://docs.aws.amazon.com/AmazonS3/latest/API/API_Owner.html
type Owner struct {
	DisplayName string
	ID          string
}

// https://docs.aws.amazon.com/AmazonS3/latest/API/ErrorResponses.html#RESTErrorResponses
type ErrorResponse struct {
	Code      string `xml:"Code"`
	Message   string `xml:"Message"`
	Resource  string `xml:"Resource"`
	RequestID string `xml:"RequestId"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_ReceiveMessage.html
type SqsReceiveMessageParams struct {
	AttributeNames              string `json:"AttributeNames,omitempty"`
	MaxNumberOfMessages         int64  `json:"MaxNumberOfMessages,omitempty"`
	MessageAttributeNames       string `json:"MessageAttributeNames,omitempty"`
	MessageSystemAttributeNames string `json:"MessageSystemAttributeNames,omitempty"`
	QueueUrl                    string `json:"QueueUrl"`
	ReceiveRequestAttemptId     string `json:"ReceiveRequestAttemptId,omitempty"`
	VisibilityTimeout           int64  `json:"VisibilityTimeout,omitempty"`
}

type MessageAttribute struct {
	// TODO: @asteurer Validate ChatGPT suggestion
	BinaryListValues [][]byte `json:"BinaryListValues"`
	BinaryValue      []byte   `json:"BinaryValue"`
	DataType         string   `json:"DataType"`
	StringListValues []string `json:"StringListValues"`
	StringValue      string   `json:"StringValue"`
}

type SqsMessage struct {
	Attributes             map[string]string           `json:"Attributes"`
	Body                   string                      `json:"Body"`
	MD5OfBody              string                      `json:"MD5OfBody"`
	MD5OfMessageAttributes string                      `json:"MD5OfMessageAttributes"`
	MessageAttributes      map[string]MessageAttribute `json:"MessageAttributes"`
	MessageId              string                      `json:"MessageId"`
	ReceiptHandle          string                      `json:"ReciptHandle"`
}

type SqsReceiveMessageResponse struct {
	Messages []SqsMessage `json:"Messages"`
}

// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_DeleteMessage.html
type SqsDeleteMessageParams struct {
	QueueUrl      string `json:"QueueUrl"`
	ReceiptHandle string `json:"ReceiptHandle"`
}

// There is not a JSON body response for a DeleteMessage request

type MessageSystemAttribute struct {
	BinaryListValues [][]byte `json:"BinaryListValues"`
	BinaryValue      []byte   `json:"BinaryValue"`
	DataType         string   `json:"DataType"`
	StringListValues []string `json:"StringListValues"`
	StringValue      string   `json:"StringValue"`
}

type SqsSendMessageParams struct {
	DelaySeconds            int64                             `json:"DelaySeconds,omitempty"`
	MessageAttributes       map[string]MessageAttribute       `json:"MessageAttributes,omitempty"`
	MessageBody             string                            `json:"MessageBody"`
	MessageDeduplicationId  string                            `json:"MessageDeduplicationId,omitempty"`
	MessageGroupId          string                            `json:"MessageGroupId,omitempty"`
	MessageSystemAttributes map[string]MessageSystemAttribute `json:"MessageSystemAttributes,omitempty"`
	QueueUrl                string                            `json:"QueueUrl"`
}

type SqsSendMessageResponse struct {
	MD5OfMessageAttributes string `json:"MD5OfMessageAttributes"`
	MD5OfMessageBody       string `json:"MD5OfMessageBody"`
	MessageId              string `json:"MessageId"`
}
