package s3

import (
	"encoding/xml"
	"reflect"
	"testing"
	"time"
)

const listBucketsResponse = `
<ListAllMyBucketsResult>
  <Buckets>
    <Bucket>
      <Name>mybucket</Name>
      <CreationDate>2024-06-24T06:34:23Z</CreationDate>
    </Bucket>
  </Buckets>
  <Owner>
    <DisplayName>webfile</DisplayName>
    <ID>75aa57f09aa0c</ID>
  </Owner>
</ListAllMyBucketsResult>
`

func TestBucketsXML(t *testing.T) {
	var info ListBucketsResponse
	if err := xml.Unmarshal([]byte(listBucketsResponse), &info); err != nil {
		t.Fatalf("failed to unmarshal: %s", err)
	}

	want := ListBucketsResponse{
		Buckets: []BucketInfo{{
			Name:         "mybucket",
			CreationDate: time.Unix(1719210863, 0).UTC(),
		}},
		Owner: Owner{
			DisplayName: "webfile",
			ID:          "75aa57f09aa0c",
		},
	}

	if !reflect.DeepEqual(info, want) {
		t.Errorf("unexpected results from parsing\ngot:  %v\nwant: %v", info, want)
	}
}

const listObjectsResponse = `
<ListBucketResult>
  <IsTruncated>true</IsTruncated>
  <Marker />
  <Name>mybucket</Name>
  <Prefix />
  <MaxKeys>1000</MaxKeys>
  <Contents>
    <Key>data.txt</Key>
    <ETag>"8ddd8be4b179a529afa5f2ffae4b9858"</ETag>
    <Owner>
      <DisplayName>webfile</DisplayName>
      <ID>75aa57f</ID>
    </Owner>
    <Size>13</Size>
    <LastModified>2024-06-24T06:34:23Z</LastModified>
    <StorageClass>STANDARD</StorageClass>
  </Contents>
  <CommonPrefixes>
    <Prefix>string</Prefix>
  </CommonPrefixes>
</ListBucketResult>
`

func TestObjectsXML(t *testing.T) {
	var info ListObjectsResponse
	if err := xml.Unmarshal([]byte(listObjectsResponse), &info); err != nil {
		t.Fatalf("failed to unmarshal: %s", err)
	}
	want := ListObjectsResponse{
		CommonPrefixes: []CommonPrefix{{Prefix: "string"}},
		Contents: []ObjectInfo{{
			Key:          "data.txt",
			ETag:         "\"8ddd8be4b179a529afa5f2ffae4b9858\"",
			Size:         13,
			LastModified: time.Unix(1719210863, 0).UTC(),
			StorageClass: "STANDARD",
			Owner: Owner{
				DisplayName: "webfile",
				ID:          "75aa57f",
			},
		}},
		IsTruncated: true,
		MaxKeys:     1000,
		Name:        "mybucket",
	}

	if !reflect.DeepEqual(info, want) {
		t.Errorf("unexpected results from parsing\ngot:  %#v\nwant: %#v", info, want)
	}
}
