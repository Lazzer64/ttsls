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

func TextDocumentDidChangeHandler(client client.Client, u protocol.Message) {
	m := u.TextDocumentDidChange()
	for _, change := range m.Params.ContentChanges.([]struct{ Text string }) {
		client.Files.Change(
			uri.URI(m.Params.TextDocument.Uri),
			change.Text,
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
			for k, category := range api.Sections {
				for _, v := range category {
					if t.Value == v.Name {
						client.Send(protocol.NewResponse(msg.Id, protocol.Hover{
							Range: protocol.Range{
								Start: protocol.Position{Character: t.Start.Character, Line: t.Start.Line},
								End:   protocol.Position{Character: t.Stop.Character + 1, Line: t.Stop.Line},
							},
							Contents: protocol.MarkupContent{
								Kind:  protocol.MarkupKindMarkdown,
								Value: v.Markdown(),
							},
						}))
						return
					}
				}
				if t.Value == k {
					names := []string{}
					for _, v := range category {
						names = append(names, v.Short())
					}
					client.Send(protocol.NewResponse(msg.Id, protocol.Hover{
						Range: protocol.Range{
							Start: protocol.Position{Character: t.Start.Character, Line: t.Start.Line},
							End:   protocol.Position{Character: t.Stop.Character + 1, Line: t.Stop.Line},
						},
						Contents: protocol.MarkupContent{
							Kind:  protocol.MarkupKindMarkdown,
							Value: fmt.Sprintf("%s\n```lua\n%s\n```", k, strings.Join(names, "\n")),
						},
					}))
					return
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

	category := api.Sections["/"]
	category = append(category, api.Sections["GlobalEvents"]...)
	if !strings.HasSuffix(string(msg.Params.TextDocument.Uri), "Global.-1.ttslua") {
		category = append(category, api.Sections["ObjectEvents"]...)
	}
	if msg.Params.Context.TriggerKind == protocol.CompletionTriggerKindTriggerCharacter {
		category = api.Sections["Object"]
	}

	for _, v := range category {
		kind := protocol.CompletionItemKindText
		switch v.Kind {
		case apiParameterKindClass:
			kind = protocol.CompletionItemKindClass
		case apiParameterKindConstant:
			kind = protocol.CompletionItemKindConstant
		case apiParameterKindEvent:
			kind = protocol.CompletionItemKindEvent
		case apiParameterKindFunction:
			kind = protocol.CompletionItemKindFunction
		case apiParameterKindProperty:
			kind = protocol.CompletionItemKindProperty
		}
		items = append(items, protocol.CompletionItem{
			Label: v.Name,
			Kind:  kind,
			Documentation: protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: v.Markdown(),
			},
		})
	}

	client.Send(protocol.NewResponse(msg.Id, items))
}

func TextDocumentSignatureHelpHandler(client client.Client, u protocol.Message) {
}

func TextDocumentCodeActionHandler(client client.Client, u protocol.Message) {
}
