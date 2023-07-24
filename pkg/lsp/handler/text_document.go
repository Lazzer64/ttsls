package handler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lazzer64/ttsls/pkg/lsp/client"
	"github.com/lazzer64/ttsls/pkg/lsp/protocol"
	"github.com/lazzer64/ttsls/pkg/lua/tokens"
	"github.com/lazzer64/ttsls/pkg/tts/script"
	"github.com/lazzer64/ttsls/pkg/uri"
)

func TextDocumentDidOpenHandler(client client.Client, u protocol.Message) {
	m := u.TextDocumentDidOpen()
	client.Files.Open(uri.URI(m.Params.TextDocument.Uri), m.Params.TextDocument.Text)
}

func TextDocumentDidCloseHandler(client client.Client, u protocol.Message) {
	m := u.TextDocumentDidClose()
	client.Files.Close(uri.URI(m.Params.TextDocument.Uri))
}

func TextDocumentDidChangeHandler(c client.Client, u protocol.Message) {
	m := u.TextDocumentDidChange()
	changes, ok := m.Params.ContentChanges.([]any)
	if !ok {
		log.Printf("LSP  Could not unmarshal content changes %#v\n", m.Params.ContentChanges)
		return
	}
	for _, change := range changes {
		c.Files.Change(
			uri.URI(m.Params.TextDocument.Uri),
			change.(map[string]any)["text"].(string),
		)
	}
}

func TextDocumentDefinitionHandler(client client.Client, u protocol.Message) {
	msg := u.TextDocumentDefinition()

	f, err := client.Files.Get(uri.URI(msg.Params.TextDocument.Uri))
	if err != nil {
		client.Log(protocol.MessageTypeError, err.Error())
		return
	}

	for _, t := range f.Tokens() {
		if t.Start.Line == msg.Params.Position.Line && t.Type == tokens.INCLUDE {
			fname := t.Value
			if !strings.HasSuffix(fname, ".ttslua") {
				fname = t.Value + ".ttslua"
			}
			if strings.HasPrefix(fname, "~/") {
				d, _ := os.UserHomeDir()
				fname = filepath.Join(d, fname[2:])
			} else {
				fname = filepath.Join(filepath.Dir(f.Path()), fname)
			}
			if _, err := os.Stat(fname); err == nil {
				log.Printf("LSP  Found source file at %s\n", fname)
				client.Send(protocol.NewResponse(msg.Id, []protocol.Location{{
					Uri: protocol.DocumentUri(fmt.Sprintf("file:///%s", fname)),
					Range: protocol.Range{
						Start: protocol.Position{Line: 0, Character: 0},
						End:   protocol.Position{Line: 0, Character: 0},
					},
				}}))
				return
			} else {
				client.Send(protocol.NewErrorResponse(
					msg.Id,
					protocol.ErrorCodesInvalidRequest,
					fmt.Errorf("Could not find file %q", fname),
				))
				return
			}
		}
	}
	client.Send(protocol.NewErrorResponse(
		msg.Id,
		protocol.LSPErrorCodesRequestFailed,
		fmt.Errorf("No definition available"),
	))
}

func TextDocumentHoverHandler(client client.Client, u protocol.Message) {
	msg := u.TextDocumentHover()

	f, err := client.Files.Get(uri.URI(msg.Params.TextDocument.Uri))
	if err != nil {
		client.Send(protocol.NewErrorResponse(
			msg.Id,
			protocol.ErrorCodesInternalError,
			fmt.Errorf("Could not read file %s: %w", msg.Params.TextDocument.Uri, err),
		))
		return
	}

	for _, t := range f.Tokens() {
		if t.Start.Line == msg.Params.Position.Line && t.Start.Character <= msg.Params.Position.Character && t.Stop.Character >= msg.Params.Position.Character {
			for c := range script.Definitions {
				for k, v := range script.Definitions[c] {
					if k == t.Value {
						client.Send(protocol.NewResponse(msg.Id, protocol.Hover{
							Range: protocol.Range{
								Start: protocol.Position{Character: t.Start.Character, Line: t.Start.Line},
								End:   protocol.Position{Character: t.Stop.Character + 1, Line: t.Stop.Line},
							},
							Contents: protocol.MarkupContent{
								Kind:  protocol.MarkupKindMarkdown,
								Value: v[0].Long,
							},
						}))
						return
					}
				}
			}
		}
	}
	client.Send(protocol.NewErrorResponse(
		msg.Id,
		protocol.LSPErrorCodesRequestFailed,
		fmt.Errorf("No hover available"),
	))
}

func TextDocumentCompletionHandler(client client.Client, u protocol.Message) {
	msg := u.TextDocumentCompletion()
	items := []protocol.CompletionItem{}
	for c := range script.Definitions {
		for _, overloads := range script.Definitions[c] {
			for _, v := range overloads {
				k := protocol.CompletionItemKindConstant
				switch v.Kind {
				case "constant":
					k = protocol.CompletionItemKindConstant
				case "property":
					k = protocol.CompletionItemKindProperty
				case "function":
					k = protocol.CompletionItemKindFunction
				case "event":
					k = protocol.CompletionItemKindEvent
				}
				items = append(items, protocol.CompletionItem{
					Label: v.Name,
					Kind:  k,
					Documentation: protocol.MarkupContent{
						Kind:  protocol.MarkupKindMarkdown,
						Value: v.Long,
					},
				})
			}
		}
	}
	client.Send(protocol.NewResponse(msg.Id, items))
}

func TextDocumentSignatureHelpHandler(c client.Client, u protocol.Message) {
	m := u.TextDocumentSignatureHelp()

	f, err := c.Files.Get(uri.URI(m.Params.TextDocument.Uri))
	if err != nil {
		c.Send(protocol.NewErrorResponse(
			m.Id,
			protocol.ErrorCodesInternalError,
			fmt.Errorf("Could not read file %s: %w", m.Params.TextDocument.Uri, err),
		))
		return
	}

	tkns := f.Tokens()
	idx := tokenIndexPre(tkns, m.Params.Position)

	// naively find the name of the function being requested by looking
	// ahead to find the number of unmatched right parens and then
	// backtracking to find the identifier preceding the closing left
	// paren backtrack through tokens to find function name.
	unclosed := 0
	for i := idx; i < len(tkns); i++ {
		if tkns[i].Type == tokens.PAREN_L {
			unclosed++
		}
		if tkns[i].Type == tokens.PAREN_R {
			unclosed--
		}
		if unclosed < 0 {
			break
		}
	}

	if unclosed == 0 {
		// no unclosed parens indicates the position is not within
		// a function.
		c.Send(protocol.NewResponse(m.Id, protocol.SignatureHelp{
			Signatures: []protocol.SignatureInformation{},
		}))
		return
	}

	var t *tokens.Token
	for i := idx; i >= 1; i-- {
		if tkns[i].Type == tokens.PAREN_L && tkns[i-1].Type == tokens.IDENTIFIER {
			t = &tkns[i-1]
			break
		}
	}

	if t == nil {
		c.Send(protocol.NewResponse(m.Id, protocol.SignatureHelp{
			Signatures: []protocol.SignatureInformation{},
		}))
		return
	}

	for category := range script.Definitions {
		for name, overloads := range script.Definitions[category] {
			if name == t.Value {
				label := name
				if category != "" {
					label = fmt.Sprintf("%s.%s", category, name)
				}
				c.Send(protocol.NewResponse(m.Id, protocol.SignatureHelp{
					ActiveSignature: 0,
					ActiveParameter: 0,
					Signatures: []protocol.SignatureInformation{{
						Label: label,
						Documentation: protocol.MarkupContent{
							Kind:  protocol.MarkupKindMarkdown,
							Value: overloads[0].Long,
						},
						Parameters:      []protocol.ParameterInformation{},
						ActiveParameter: 0,
					}},
				}))
			}
		}
	}

	c.Send(protocol.NewResponse(m.Id, protocol.SignatureHelp{
		Signatures: []protocol.SignatureInformation{},
	}))
}

func TextDocumentCodeActionHandler(client client.Client, u protocol.Message) {
}

// tokenIndex returns the index of the token at position p, or -1 if none exists
func tokenIndex(tkns []tokens.Token, p protocol.Position) int {
	for i, t := range tkns {
		if t.Start.Line == p.Line && t.Start.Character <= p.Character && t.Stop.Character >= p.Character {
			return i
		}
	}
	return -1
}

// tokenIndexPre returns the index of the token at or preceding position p
func tokenIndexPre(tkns []tokens.Token, p protocol.Position) int {
	for i := 1; i < len(tkns); i++ {
		nxt := tkns[i]
		if nxt.Start.Line > p.Line || (nxt.Start.Line == p.Line && nxt.Start.Character > p.Character) {
			return i-1
		}
	}
	return len(tkns)-1
}
