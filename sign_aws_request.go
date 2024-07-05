package aws

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

func GetAuthorizationHeader(config *Config, req *http.Request, date *AwsDate, payloadHash string) string {
	canonicalHeaders, signedHeaders := getHeaderStrings(req.Header)
	canonicalRequest := getCanonicalRequest(req, signedHeaders, canonicalHeaders, payloadHash)
	stringToSign := getStringToSign(config, date, canonicalRequest)
	signature := getSignature(config, date, stringToSign)
	credential := strings.Join([]string{config.AccessKeyId, date.GetDate(), config.Region, config.Service, "aws4_request"}, "/")

	return fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s, SignedHeaders=%s, Signature=%s",
		credential, signedHeaders, signature)
}

func getHeaderStrings(headers http.Header) (string, string) {
	// Formatted as header_key_1:header_value_1\nheader_key_2:header_value_2\n
	canonicalHeaders := ""
	// Formatted as header_key_1;header_key_2
	signedHeaders := ""
	headerKeys := make([]string, 0, len(headers))
	for key := range headers {
		headerKeys = append(headerKeys, key)
	}
	// Header names must appear in alphabetical order
	sort.Strings(headerKeys)

	for _, key := range headerKeys {
		// Each header name must use lowercase characters
		lowerCaseKey := strings.ToLower(key)
		canonicalHeaders += lowerCaseKey + ":" + headers.Get(key) + "\n"
		if signedHeaders == "" {
			signedHeaders += lowerCaseKey
		} else {
			signedHeaders += ";" + lowerCaseKey
		}
	}

	return canonicalHeaders, signedHeaders
}

// https://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-header-based-auth.html#request-string
func getStringToSign(config *Config, date *AwsDate, canonicalRequest string) string {
	scope := strings.Join([]string{date.GetDate(), config.Region, config.Service, "aws4_request"}, "/")
	return strings.Join([]string{"AWS4-HMAC-SHA256", date.GetTime(), scope, getPayloadHash([]byte(canonicalRequest))}, "\n")
}

// https://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-header-based-auth.html#signing-key
func getSignature(config *Config, date *AwsDate, stringToSign string) string {
	sign := func(key []byte, data []byte) []byte {
		hash := hmac.New(sha256.New, key)
		hash.Write(data)

		return hash.Sum(nil)
	}

	dateKey := sign([]byte("AWS4"+config.SecretAccessKey), []byte(date.GetDate()))
	regionKey := sign(dateKey, []byte(config.Region))
	serviceKey := sign(regionKey, []byte(config.Service))
	signingKey := sign(serviceKey, []byte("aws4_request"))

	return hex.EncodeToString(sign(signingKey, []byte(stringToSign)))
}

// https://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-header-based-auth.html#canonical-request
func getCanonicalRequest(req *http.Request, signedHeaders, canonicalHeaders, payloadHash string) string {
	escapedUrl := req.URL.EscapedPath()
	if !strings.HasPrefix(escapedUrl, "/") {
		// The path MUST start with a "/"
		escapedUrl = "/" + escapedUrl
	}

	return strings.Join([]string{
		req.Method,
		escapedUrl,
		req.URL.RawQuery,
		canonicalHeaders,
		signedHeaders,
		payloadHash,
	}, "\n")
}

func getPayloadHash(payload []byte) string {
	hash := sha256.New()
	hash.Write(payload)
	return hex.EncodeToString(hash.Sum(nil))
}
