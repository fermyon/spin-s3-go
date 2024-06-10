package s3

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func localClient(t *testing.T) *Client {
	t.Helper()

	cfg := Config{
		Endpoint: "http://s3.localhost.localstack.cloud:4566",
	}
	client, err := New(cfg)
	require.NoError(t, err, "failed creating test client")

	return client
}

func awsClient(t *testing.T) *Client {
	t.Helper()

	// aws config
	cfg := Config{
		AccessKey:    os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey:    os.Getenv("AWS_SECRET_ACCESS_KEY"),
		SessionToken: os.Getenv("AWS_SESSION_TOKEN"),
		Region:       "us-east-1",
		// Endpoint:     "https://s3.dualstack.us-east-1.amazonaws.com",

		Endpoint: "https://s3.us-east-1.amazonaws.com",
	}
	client, err := New(cfg)
	require.NoError(t, err, "failed creating test client")

	return client
}

func TestListBuckets(t *testing.T) {
	t.Skip("tests will only work if spinhttp is replaced with default http")
	client := localClient(t)

	ctx := context.Background()
	err := client.CreateBucket(ctx, "foo")
	require.NoError(t, err, "failed listing buckets")

	err = client.CreateBucket(ctx, "bar")
	require.NoError(t, err, "failed listing buckets")

	// Trace the next request
	client.trace = true
	resp, err := client.ListBuckets(ctx)
	require.NoError(t, err, "failed listing buckets")
	require.NotNil(t, resp)
	// TODO assert the response

	t.Log("List buckets response:")
	t.Logf("%#v", resp.Buckets)

	// TODO cleanup buckets
}

func TestListObjects(t *testing.T) {
	t.Skip("tests will only work if spinhttp is replaced with default http")
	client := localClient(t)

	ctx := context.Background()
	err := client.CreateBucket(ctx, "foo")
	require.NoError(t, err, "failed listing buckets")

	err = client.PutObject(ctx, "foo", "data.txt", strings.NewReader("Hello S3!"))
	require.NoError(t, err, "failed putting object")

	// Trace the next request
	client.trace = true
	listObjectsResp, err := client.ListObjects(ctx, "foo")
	require.NoError(t, err, "failed listing objects")
	require.NotNil(t, listObjectsResp)
	// TODO assert the response

	t.Log("List objects response:")
	t.Logf("%#v", listObjectsResp)

	// TODO cleanup buckets
}

func TestClientLocal(t *testing.T) {
	t.Skip("tests will only work if spinhttp is replaced with default http")
	client := localClient(t)

	ctx := context.Background()
	err := client.CreateBucket(ctx, "foo")
	require.NoError(t, err, "failed listing buckets")

	r := strings.NewReader("Hello S3!")
	err = client.PutObject(ctx, "foo", "data.txt", r)
	require.NoError(t, err, "failed putting object")

	object, err := client.GetObject(ctx, "foo", "data.txt")
	require.NoError(t, err, "failed getting object")

	b, err := io.ReadAll(object)
	require.NoError(t, err, "failed reading object")
	t.Logf("%#v", string(b))

	listObjectsResp, err := client.ListObjects(ctx, "foo")
	require.NoError(t, err, "failed listing objects")
	require.NotNil(t, listObjectsResp)

	t.Log("List objects response:")
	t.Logf("%#v", listObjectsResp)

	// TODO cleanup buckets
}

func TestClientAWS(t *testing.T) {
	t.Skip("tests will only work if spinhttp is replaced with default http")

	client := awsClient(t)
	client.trace = true

	// err = client.CreateBucket(context.TODO(), "foo")
	// require.NoError(t, err, "failed listing buckets")

	// err = client.CreateBucket(context.TODO(), "bar")
	// require.NoError(t, err, "failed listing buckets")

	buckets, err := client.ListBuckets(context.TODO())
	require.NoError(t, err, "failed listing buckets")
	require.NotNil(t, buckets)

	t.Logf("%#v", buckets)
}
