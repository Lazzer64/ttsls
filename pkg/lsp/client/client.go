package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/lazzer64/ttsls/pkg/lsp/protocol"
	"github.com/lazzer64/ttsls/pkg/uri"
)

type Client struct {
	w     io.Writer
	Files Files
}

func New(w io.Writer) Client {
	return Client{w, &fs{w, map[uri.URI]SourceFile{}}}
}

func (c Client) Send(v any) error {
	return writeMsg(c.w, v)
}

func (c Client) Log(level protocol.MessageType, msg string) error {
	return c.Send(
		protocol.NewWindowShowMessageRequestRequest(
			0,
			protocol.ShowMessageRequestParams{
				Type:    level,
				Message: msg,
			},
		),
	)
}

func writeMsg(w io.Writer, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		log.Printf("Could not marshal data: %s", err)
		return err
	}
	_, err = fmt.Fprintf(w, fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(b), b))
	if err != nil {
		log.Printf("Could not write response: %s", err)
		return err
	}
	log.Printf("SEND %s\n", b)
	return nil
}
