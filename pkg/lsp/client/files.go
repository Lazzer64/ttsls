package client

import (
	"io"
	"math/rand"
	"os"
	"strings"
	"sync"

	"github.com/lazzer64/ttsls/pkg/lsp/protocol"
	"github.com/lazzer64/ttsls/pkg/lua/tokens"
	"github.com/lazzer64/ttsls/pkg/uri"
)

type Files interface {
	Open(uri.URI, string)
	Close(uri.URI)
	Change(uri.URI, string)
	Get(uri.URI) (SourceFile, error)
	GetAllOpen() []SourceFile
	Write(uri.URI, string) error
}

type fs struct {
	w       io.Writer
	managed map[uri.URI]SourceFile
	mtx     sync.Mutex
}

func (w *fs) Open(uri uri.URI, text string) {
	w.mtx.Lock()
	w.managed[uri] = SourceFile{
		uri:    uri,
		text:   text,
		tokens: tokens.TokenizeString(text),
	}
	w.mtx.Unlock()
}

func (w *fs) Close(uri uri.URI) {
	delete(w.managed, uri)
}

func (w *fs) Change(uri uri.URI, text string) {
	w.mtx.Lock()
	w.managed[uri] = SourceFile{
		uri:    uri,
		text:   text,
		tokens: tokens.TokenizeString(text),
	}
	w.mtx.Unlock()
}

func (w *fs) Write(uri uri.URI, text string) error {
	if f, ok := w.managed[uri]; ok {
		return writeMsg(w.w, protocol.NewWorkspaceApplyEditRequest(
			rand.Int(), // TODO: replace rand.Int
			protocol.ApplyWorkspaceEditParams{
				Edit: protocol.WorkspaceEdit{
					Changes: []protocol.TextEdit{{
						Range:   protocol.Range{Start: f.Start(), End: f.End()},
						NewText: text,
					}},
				},
			},
		))
	}
	return os.WriteFile(uri.Path(), []byte(text), 0644)
}

func (w *fs) Get(uri uri.URI) (SourceFile, error) {
	if sf, ok := w.managed[uri]; ok {
		return sf, nil
	}
	b, err := os.ReadFile(uri.Path())
	if err != nil {
		return SourceFile{}, err
	}
	return SourceFile{
		uri:    uri,
		text:   string(b),
		tokens: tokens.Tokenize(b),
	}, nil
}

func (w *fs) GetAllOpen() []SourceFile {
	files := make([]SourceFile, 0, len(w.managed))
	for _, v := range w.managed {
		files = append(files, v)
	}
	return files
}

type SourceFile struct {
	uri    uri.URI
	text   string
	tokens []tokens.Token
}

func ReadSourceFile(name string) (SourceFile, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return SourceFile{}, err
	}
	uri, err := uri.Parse(name)
	if err != nil {
		return SourceFile{}, err
	}
	return SourceFile{
		uri:    uri,
		text:   string(b),
		tokens: tokens.Tokenize(b),
	}, err
}

func (f *SourceFile) URI() uri.URI {
	return f.uri
}

func (f *SourceFile) Path() string {
	return f.uri.Path()
}

func (f *SourceFile) Content() string {
	return f.text
}

func (f *SourceFile) Tokens() []tokens.Token {
	return f.tokens
}

func (f *SourceFile) Start() protocol.Position {
	return protocol.Position{Line: 0, Character: 0}
}

func (f *SourceFile) End() protocol.Position {
	lines := strings.Split(f.text, "\n")
	return protocol.Position{Line: uint32(len(lines) + 1), Character: 0}
}
