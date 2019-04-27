package mdata

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/tinylib/msgp/msgp"

	"github.com/raintank/schema"
)

type ChunkSaveCallback func(error)

// ChunkWriteRequest is a request to write a chunk into a store
//go:generate msgp
type ChunkWriteRequest struct {
	Callback  ChunkSaveCallback `msg:"-"`
	Key       schema.AMKey
	TTL       uint32
	T0        uint32
	Data      []byte
	Timestamp time.Time
}

// NewChunkWriteRequest creates a new ChunkWriteRequest
func NewChunkWriteRequest(callback ChunkSaveCallback, key schema.AMKey, ttl, t0 uint32, data []byte, ts time.Time) ChunkWriteRequest {
	return ChunkWriteRequest{callback, key, ttl, t0, data, ts}
}

// ArchiveRequest is used by the whisper importer utilities to transfer the data
// from the reader to the writer
//go:generate msgp
type ArchiveRequest struct {
	MetricData         schema.MetricData
	ChunkWriteRequests []ChunkWriteRequest
}

func (a *ArchiveRequest) MarshalCompressed() (*bytes.Buffer, error) {
	var buf bytes.Buffer

	buf.WriteByte(byte(uint8(1)))

	g := gzip.NewWriter(&buf)
	msgp.Encode(g, a)

	err := g.Close()
	if err != nil {
		return &buf, fmt.Errorf("ERROR: Compressing MSGP data: %q", err)
	}

	return &buf, nil
}

func (a *ArchiveRequest) UnmarshalCompressed(b io.Reader) error {
	versionBuf := make([]byte, 1)
	readBytes, err := b.Read(versionBuf)
	if err != nil || readBytes != 1 {
		return fmt.Errorf("ERROR: Failed to read one byte: %s", err)
	}

	version := uint8(versionBuf[0])
	if version != 1 {
		return errors.New(fmt.Sprintf("ERROR: Only version 1 is supported, received version %d", version))
	}

	gzipReader, err := gzip.NewReader(b)
	if err != nil {
		return fmt.Errorf("ERROR: Creating Gzip reader: %q", err)
	}

	err = msgp.Decode(bufio.NewReader(gzipReader), a)
	if err != nil {
		return fmt.Errorf("ERROR: Unmarshaling Raw: %q", err)
	}
	gzipReader.Close()

	return nil
}
