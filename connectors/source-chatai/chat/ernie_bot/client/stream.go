// Copyright 2022 Linkall Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ChatCompletionStream struct {
	isFinished     bool
	reader         *bufio.Reader
	response       *http.Response
	errAccumulator errorAccumulator
}

func (stream *ChatCompletionStream) Recv() (*ChatCompletionResponse, error) {
	if stream.isFinished {
		return nil, io.EOF
	}
	var headerData = []byte("data: ")

waitForData:
	line, err := stream.reader.ReadBytes('\n')
	if err != nil {
		if writeErr := stream.errAccumulator.write(line); writeErr != nil {
			return nil, writeErr
		}
		respErr := stream.errAccumulator.unmarshalError()
		if respErr != nil {
			err = fmt.Errorf("response error code:%d, msg:%s", respErr.ErrorCode, respErr.ErrorMsg)
		}
		return nil, err
	}
	line = bytes.TrimSpace(line)
	if !bytes.HasPrefix(line, headerData) {
		if writeErr := stream.errAccumulator.write(line); writeErr != nil {
			return nil, writeErr
		}
		goto waitForData
	}
	line = bytes.TrimPrefix(line, headerData)
	var response ChatCompletionResponse
	if err = json.Unmarshal(line, &response); err != nil {
		return nil, err
	}
	if response.IsEnd {
		stream.isFinished = true
	}
	return &response, nil
}

func (stream *ChatCompletionStream) Close() {
	stream.response.Body.Close()
}
