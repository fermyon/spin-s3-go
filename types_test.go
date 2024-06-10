package s3

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/require"
)

const listBucketsResponse = `
<ListAllMyBucketsResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
  <Owner>
    <DisplayName>webfile</DisplayName>
    <ID>75aa57f09aa0c8caeab4f8c24e99d10f8e7faeebf76c078efc7c6caea54ba06a</ID>
  </Owner>
  <Buckets>
    <Bucket>
      <Name>mybucket</Name>
      <CreationDate>2024-05-13T20:42:12.000Z</CreationDate>
    </Bucket>
    <Bucket>
      <Name>foo</Name>
      <CreationDate>2024-05-13T20:58:39.000Z</CreationDate>
    </Bucket>
    <Bucket>
      <Name>bar</Name>
      <CreationDate>2024-05-13T21:48:41.000Z</CreationDate>
    </Bucket>
  </Buckets>
</ListAllMyBucketsResult>
`

const listObjectsResponse = `
<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
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
      <ID>
      75aa57f09aa0c8caeab4f8c24e99d10f8e7faeebf76c078efc7c6caea54ba06a</ID>
    </Owner>
    <Size>13</Size>
    <LastModified>2024-05-13T20:46:46.000Z</LastModified>
    <StorageClass>STANDARD</StorageClass>
  </Contents>
  <CommonPrefixes>
    <Prefix>string</Prefix>
  </CommonPrefixes>
</ListBucketResult>
`

func TestXML(t *testing.T) {
	t.Skip()
	var info ListBucketsResponse
	err := xml.Unmarshal([]byte(listBucketsResponse), &info)
	require.NoError(t, err, "Unmarshal failed")
	t.Logf("%#v", info)

	out, err := xml.MarshalIndent(&info, "  ", "    ")
	require.NoError(t, err, "Marshal failed")
	t.Log(string(out))
}

func TestBucketsXML(t *testing.T) {
	t.Skip()
	var info ListObjectsResponse
	err := xml.Unmarshal([]byte(listBucketsResponse), &info)
	require.NoError(t, err, "Unmarshal failed")
	t.Logf("%#v", info)

	// info.CommonPrefixes = []CommonPrefix{{Prefix: "hello"}}

	out, err := xml.MarshalIndent(&info, "  ", "    ")
	require.NoError(t, err, "Marshal failed")
	t.Log(string(out))
}
