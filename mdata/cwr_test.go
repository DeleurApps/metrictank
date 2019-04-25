package mdata

import (
	"bytes"
	"testing"

	"github.com/raintank/schema"
)

func TestArchiveRequestEncodingDecoding(t *testing.T) {
	originalRequest := ArchiveRequest{
		MetricData: schema.MetricData{
			Name: "testMetricData",
		},
		ChunkWriteRequests: []ChunkWriteRequest{
			{
				TTL: 1,
				T0:  2,
			},
		},
	}

	encoded, err := originalRequest.MarshalCompressed()
	reader := bytes.NewReader(encoded.Bytes())
	if err != nil {
		t.Fatalf("Expected no error when encoding request, got %q", err)
	}

	decodedRequest := ArchiveRequest{}
	err = decodedRequest.UnmarshalCompressed(reader)
	if err != nil {
		t.Fatalf("Expected no error when decoding request, got %q", err)
	}

	if originalRequest.MetricData.Name != decodedRequest.MetricData.Name || originalRequest.ChunkWriteRequests[0].TTL != decodedRequest.ChunkWriteRequests[0].TTL || originalRequest.ChunkWriteRequests[0].T0 != decodedRequest.ChunkWriteRequests[0].T0 {
		t.Fatalf("Decoded request is different than the encoded one")
	}
}

func TestArchiveRequestEncodingWithUnknownVersion(t *testing.T) {
	originalRequest := ArchiveRequest{}
	encoded, err := originalRequest.MarshalCompressed()
	if err != nil {
		t.Fatalf("Expected no error when encoding request, got %q", err)
	}

	encodedBytes := encoded.Bytes()
	encodedBytes[0] = byte(uint8(2))
	decodedRequest := ArchiveRequest{}
	err = decodedRequest.UnmarshalCompressed(bytes.NewReader(encodedBytes))
	if err == nil {
		t.Fatalf("Expected an error when decoding, but got nil")
	}
}
