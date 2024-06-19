package s3

import (
	"context"
	"io"
	"os"
	"testing"
)

func localClient(t *testing.T) *Client {
	t.Helper()

	cfg := Config{
		Endpoint: "http://s3.localhost.localstack.cloud:4566",
	}
	client, err := New(cfg)
	if err != nil {
		t.Fatalf("failed creating test client: %s", err)
	}
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
	if err != nil {
		t.Fatalf("failed creating test client: %s", err)
	}

	return client
}

func TestListBuckets(t *testing.T) {
	t.Skip("tests will only work if spinhttp is replaced with default http")
	client := localClient(t)

	ctx := context.Background()
	if err := client.CreateBucket(ctx, "foo"); err != nil {
		t.Fatalf("failed creating bucket: %s", err)
	}

	if err := client.CreateBucket(ctx, "bar"); err != nil {
		t.Fatalf("failed creating bucket: %s", err)
	}

	// Trace the next request
	client.trace = true
	resp, err := client.ListBuckets(ctx)
	if err != nil {
		t.Fatalf("failed listing buckets: %s", err)
	}
	// TODO assert the response

	t.Log("List buckets response:")
	t.Logf("%#v", resp.Buckets)

	// TODO cleanup buckets
}

func TestListObjects(t *testing.T) {
	t.Skip("tests will only work if spinhttp is replaced with default http")
	client := localClient(t)

	ctx := context.Background()
	if err := client.CreateBucket(ctx, "foo"); err != nil {
		t.Fatalf("failed creating bucket: %s", err)
	}

	if err := client.PutObject(ctx, "foo", "data.txt", []byte("Hello S3!")); err != nil {
		t.Fatalf("failed putting object: %s", err)
	}

	// Trace the next request
	client.trace = true
	listObjectsResp, err := client.ListObjects(ctx, "foo")
	if err != nil {
		t.Fatalf("failed listing objects: %s", err)
	}
	// TODO assert the response

	t.Log("List objects response:")
	t.Logf("%#v", listObjectsResp)

	// TODO cleanup buckets
}

func TestClientLocal(t *testing.T) {
	t.Skip("tests will only work if spinhttp is replaced with default http")
	client := localClient(t)

	ctx := context.Background()
	if err := client.CreateBucket(ctx, "foo"); err != nil {
		t.Fatalf("failed creating bucket: %s", err)
	}

	if err := client.PutObject(ctx, "foo", "data.txt", []byte("Hello S3!")); err != nil {
		t.Fatalf("failed putting object: %s", err)
	}

	object, err := client.GetObject(ctx, "foo", "data.txt")
	if err != nil {
		t.Fatalf("failed getting object: %s", err)
	}

	b, err := io.ReadAll(object)
	if err != nil {
		t.Fatalf("failed reading object: %s", err)
	}
	t.Logf("%#v", string(b))

	listObjectsResp, err := client.ListObjects(ctx, "foo")
	if err != nil {
		t.Fatalf("failed listing objects: %s", err)
	}

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
	if err != nil {
		t.Fatalf("failed listing buckets: %s", err)
	}

	t.Logf("%#v", buckets)
}
