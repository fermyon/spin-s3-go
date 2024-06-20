package s3

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
	FileSize     int
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
