package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/lazzer64/ttsls/pkg/lsp/protocol"
)

type MessageWriter interface {
	Send(v any) error
	Log(level protocol.MessageType, msg string) error
}

type messageWriter struct {
	io.Writer
}

func (w messageWriter) Send(v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		log.Printf("Could not marshal data: %s", err)
		return err
	}
	_, err = w.Write([]byte(fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(b), b)))
	if err != nil {
		log.Printf("Could not write response: %s", err)
		return err
	}
	log.Printf("SEND %s\n", b)
	return nil
}

func (w messageWriter) Log(level protocol.MessageType, msg string) error {
	return w.Send(
		protocol.NewWindowShowMessageRequestRequest(
			0,
			protocol.ShowMessageRequestParams{
				Type:    level,
				Message: msg,
			},
		),
	)
}
