package mdata

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/raintank/schema"
)

type ChunkSaveCallback func()

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
func NewChunkWriteRequest(callback func(), key schema.AMKey, ttl, t0 uint32, data []byte, ts time.Time) ChunkWriteRequest {
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

	// prefix the buffer with version number 1
	buf.WriteByte(byte(uint8(1)))

	b, err := a.MarshalMsg(nil)
	if err != nil {
		return &buf, errors.New(fmt.Sprintf("ERROR: Marshalling metric: %q", err))
	}

	g := gzip.NewWriter(&buf)
	_, err = g.Write(b)
	if err != nil {
		return &buf, errors.New(fmt.Sprintf("ERROR: Compressing MSGP data: %q", err))
	}

	err = g.Close()
	if err != nil {
		return &buf, errors.New(fmt.Sprintf("ERROR: Compressing MSGP data: %q", err))
	}

	return &buf, nil
}

func (a *ArchiveRequest) UnmarshalCompressed(b io.Reader) error {
	reader := bufio.NewReader(b)
	versionBuf, err := reader.ReadByte()
	if err != nil {
		return errors.New(fmt.Sprintf("ERROR: Failed to read received data: %q", err))
	}

	version := uint8(versionBuf)
	if version != 1 {
		return errors.New(fmt.Sprintf("ERROR: Only version 1 is supported, received version %d", version))
	}

	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return errors.New(fmt.Sprintf("ERROR: Creating Gzip reader: %q", err))
	}

	raw, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return errors.New(fmt.Sprintf("ERROR: Decompressing Gzip: %q", err))
	}

	_, err = a.UnmarshalMsg(raw)
	if err != nil {
		return errors.New(fmt.Sprintf("ERROR: Unmarshaling Raw: %q", err))
	}

	return nil
}
