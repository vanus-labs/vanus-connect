package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type errorAccumulator interface {
	write(p []byte) error
	unmarshalError() *ErrResponse
}

type errorBuffer interface {
	io.Writer
	Len() int
	Bytes() []byte
}

type defaultErrorAccumulator struct {
	buffer errorBuffer
}

func newErrorAccumulator() errorAccumulator {
	return &defaultErrorAccumulator{
		buffer: &bytes.Buffer{},
	}
}

func (e *defaultErrorAccumulator) write(p []byte) error {
	_, err := e.buffer.Write(p)
	if err != nil {
		return fmt.Errorf("error accumulator write error, %w", err)
	}
	return nil
}

func (e *defaultErrorAccumulator) unmarshalError() (errResp *ErrResponse) {
	if e.buffer.Len() == 0 {
		return
	}
	err := json.Unmarshal(e.buffer.Bytes(), &errResp)
	if err != nil {
		errResp = nil
	}

	return
}
