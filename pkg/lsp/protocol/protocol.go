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

type Request struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	raw     []byte
}

func ReadRequest(r io.Reader) (Request, error) {
	length := 0
	resp := Request{}

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
