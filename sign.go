package s3

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

// Lifted from https://github.com/joshuarose/spin-go-aws/blob/main/aws-no-sdk/main.go
// It works for all requests but GetObject() ಠ_ಠ

func (c *Client) getAuthorizationHeader(req *http.Request, now time.Time) string {
	// Create the canonical request
	canonicalRequest := getCanonicalRequest(req.Host, req.Method, c.config.SessionToken, now)

	// Create the string to sign
	stringToSign := getStringToSign(canonicalRequest, c.config.Region, now)

	// Calculate the signature
	signature := getSignature(stringToSign, c.config.Region, c.config.SecretKey, now)

	// Create the authorization header
	authorizationHeader := fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s/%s/s3/aws4_request, SignedHeaders=host;x-amz-content-sha256;x-amz-date;x-amz-security-token, Signature=%s",
		c.config.AccessKey, now.Format("20060102"), c.config.Region, signature)

	return authorizationHeader
}

func getStringToSign(canonicalRequest, region string, now time.Time) string {
	// Create the hash of the canonical request
	canonicalRequestHash := sha256.New()
	canonicalRequestHash.Write([]byte(canonicalRequest))
	canonicalRequestHashString := hex.EncodeToString(canonicalRequestHash.Sum(nil))

	// Create the string to sign
	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s/%s/s3/aws4_request\n%s",
		now.Format("20060102T150405Z"), now.Format("20060102"), region, canonicalRequestHashString)

	return stringToSign
}

func getSignature(stringToSign, region, secretKey string, now time.Time) string {
	// Create the signing key
	dateKey := hmacSHA256([]byte("AWS4"+secretKey), []byte(now.Format("20060102")))
	regionKey := hmacSHA256(dateKey, []byte(region))
	serviceKey := hmacSHA256(regionKey, []byte("s3"))
	signingKey := hmacSHA256(serviceKey, []byte("aws4_request"))

	// Calculate the signature
	signature := hmacSHA256(signingKey, []byte(stringToSign))

	return hex.EncodeToString(signature)
}

func getCanonicalRequest(host, method, sessionToken string, now time.Time) string {
	// Create the canonical URI
	canonicalURI := "/"

	// Create the canonical query string
	canonicalQueryString := ""

	// Create the canonical headers
	canonicalHeaders := fmt.Sprintf("host:%s\nx-amz-content-sha256:%s\nx-amz-date:%s\nx-amz-security-token:%s\n",
		host, getPayloadHash(""), now.Format("20060102T150405Z"), sessionToken)

	// Create the signed headers
	signedHeaders := "host;x-amz-content-sha256;x-amz-date;x-amz-security-token"

	// Create the payload hash
	payloadHash := sha256.New()
	payloadHash.Write([]byte(""))
	payloadHashString := hex.EncodeToString(payloadHash.Sum(nil))

	// Combine all the components to create the canonical request
	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		method, canonicalURI, canonicalQueryString, canonicalHeaders, signedHeaders, payloadHashString)

	return canonicalRequest
}

func hmacSHA256(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

func getPayloadHash(payload string) string {
	hash := sha256.New()
	hash.Write([]byte(payload))
	return hex.EncodeToString(hash.Sum(nil))
}
