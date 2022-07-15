package client

import (
	"io"

	"github.com/lazzer64/ttsls/pkg/uri"
)

type Client struct {
	MessageWriter
	Files Files
}

func New(w io.Writer) Client {
	mw := messageWriter{w}
	return Client{mw, &fs{mw, map[uri.URI]SourceFile{}}}
}
