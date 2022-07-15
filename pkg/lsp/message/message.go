package message

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type Message interface {
	Bytes() []byte
}

type RequestMessage struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

func (m RequestMessage) Bytes() []byte {
	b, _ := json.Marshal(m)
	return b
}

func NewRequestMessage(id int, method string, params interface{}) RequestMessage {
	return RequestMessage{
		Jsonrpc: "2.0",
		Id:      id,
		Method:  method,
		Params:  params,
	}
}

type ResponseMessage struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      int         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func NewResponseMessage(id int, result interface{}, err interface{}) ResponseMessage {
	return ResponseMessage{
		Jsonrpc: "2.0",
		Id:      id,
		Result:  result,
		Error:   err,
	}
}

func (m ResponseMessage) Bytes() []byte {
	b, _ := json.Marshal(m)
	return b
}

type NotificationMessage struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

func (m NotificationMessage) Bytes() []byte {
	b, _ := json.Marshal(m)
	return b
}

func NewNotificationMessage(method string, params interface{}) NotificationMessage {
	return NotificationMessage{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
	}
}

type ShowMessageLevel int

const (
	MessageTypeError   ShowMessageLevel = 1
	MessageTypeWarning ShowMessageLevel = 2
	MessageTypeInfo    ShowMessageLevel = 3
	MessageTypeLog     ShowMessageLevel = 4
)

func NewShowMessageNotification(level ShowMessageLevel, msg string) NotificationMessage {
	params := map[string]interface{}{
		"type":    level,
		"message": msg,
	}
	return NewNotificationMessage("window/showMessage", params)
}

type UndefinedMessage struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	raw     []byte
}

func (m UndefinedMessage) Bytes() []byte {
	return m.raw
}

func ReadMessage(r io.Reader) (UndefinedMessage, error) {
	length := 0
	resp := UndefinedMessage{}

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
