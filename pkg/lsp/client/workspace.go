package client

import (
	"math/rand"
	"os"
	"strings"

	"github.com/lazzer64/ttsls/pkg/lsp/message"
	"github.com/lazzer64/ttsls/pkg/lsp/message/types"
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
	mw      MessageWriter
	managed map[uri.URI]SourceFile
}

func (w *fs) Open(uri uri.URI, text string) {
	w.managed[uri] = SourceFile{
		uri:    uri,
		text:   text,
		tokens: tokens.TokenizeString(text),
	}
}

func (w *fs) Close(uri uri.URI) {
	delete(w.managed, uri)
}

func (w *fs) Change(uri uri.URI, text string) {
	w.managed[uri] = SourceFile{
		uri:    uri,
		text:   text,
		tokens: tokens.TokenizeString(text),
	}
}

func (w *fs) Write(uri uri.URI, text string) error {
	if f, ok := w.managed[uri]; ok {
		w.mw.Send(message.NewRequestMessage(
			rand.Int(), // TODO: replace rand.Int
			"workspace/applyEdit",
			map[string]any{
				"edit": map[string]any{
					"changes": []types.TextEdit{{
						Range: types.Range{
							Start: f.Start(),
							End:   f.End(),
						},
						NewText: text,
					}},
				},
			},
		))
		return nil
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

func (f *SourceFile) Start() types.Position {
	return types.Position{Line: 0, Character: 0}
}

func (f *SourceFile) End() types.Position {
	lines := strings.Split(f.text, "\n")
	return types.Position{Line: len(lines) + 1, Character: 0}
}
