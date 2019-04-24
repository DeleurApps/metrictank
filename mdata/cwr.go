package mdata

import (
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

// MdWithCwrs is used by the whisper importer utilities to transfer the data
// from the reader to the writer
//go:generate msgp
type MdWithCwrs struct {
	Md   schema.MetricData
	Cwrs []ChunkWriteRequest
}

func (m *MdWithCwrs) MarshalCompressed() (*bytes.Buffer, error) {
	var buf bytes.Buffer

	b, err := m.MarshalMsg(nil)
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

func (m *MdWithCwrs) UnmarshalCompressed(b io.Reader) error {
	gzipReader, err := gzip.NewReader(b)
	if err != nil {
		return errors.New(fmt.Sprintf("ERROR: Creating Gzip reader: %q", err))
	}

	raw, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return errors.New(fmt.Sprintf("ERROR: Decompressing Gzip: %q", err))
	}

	_, err = m.UnmarshalMsg(raw)
	if err != nil {
		return errors.New(fmt.Sprintf("ERROR: Unmarshaling Raw: %q", err))
	}

	return nil
}
