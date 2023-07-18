package protocol

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

//go:generate go run generate/main.go release/protocol/3.17.3

type DocumentUri string
type URI string

type Message struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	raw     []byte
}

func ReadMessage(r io.Reader) (Message, error) {
	length := 0
	resp := Message{}

	_, err := fmt.Fscanf(r, "Content-Length: %d\r\n\r\n", &length)
	// TODO: handle optional Content-Type

	if err != nil {
		return resp, errors.New(fmt.Sprintf("Could not determine content length: %s\n", err))
	}

	buf := make([]byte, length)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return resp, errors.New(fmt.Sprintf("Expected %d bytes but read %d: %s\n", length, n, err))
	}

	json.Unmarshal(buf, &resp)
	resp.raw = buf
	return resp, nil
}

type Response struct {
	Id      int            `json:"id"`
	Jsonrpc string         `json:"jsonrpc"`
	Result  any            `json:"result,omitempty"`
	Error   *ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Code    int `json:"code"`
	Message string     `json:"message"`
	Data    any        `json:"data,omitempty"`
}

func NewResponse(id int, result any) Response {
	return Response{
		Id:      id,
		Jsonrpc: "2.0",
		Result:  result,
		Error:   nil,
	}
}

func NewErrorResponse[T ~int](id int, code T, err error) Response {
	return Response{
		Id:      id,
		Jsonrpc: "2.0",
		Result:  nil,
		Error: &ResponseError{
			Code:    int(code),
			Message: err.Error(),
			Data:    err.Error(),
		},
	}
}
