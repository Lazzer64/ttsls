package handler

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lazzer64/ttsls/pkg/lsp/client"
	"github.com/lazzer64/ttsls/pkg/lsp/message"
	"github.com/lazzer64/ttsls/pkg/lsp/message/types"
	"github.com/lazzer64/ttsls/pkg/lua/tokens"
)

func TextDocumentDidOpenHandler(client client.Client, u message.UndefinedMessage) {
	m := u.TextDocumentDidOpenMessage()
	client.Files.Open(m.Params.TextDocument.Uri, m.Params.TextDocument.Text)
}

func TextDocumentDidCloseHandler(client client.Client, u message.UndefinedMessage) {
	m := u.TextDocumentDidCloseMessage()
	client.Files.Close(m.Params.TextDocument.Uri)
}

func TextDocumentDidChangeHandler(client client.Client, u message.UndefinedMessage) {
	m := u.TextDocumentDidChangeMessage()
	for _, change := range m.Params.ContentChanges {
		client.Files.Change(
			m.Params.TextDocument.Uri,
			change.Text,
		)
	}
}

func TextDocumentDefinitionHandler(client client.Client, u message.UndefinedMessage) {
	msg := u.TextDocumentDefinition()

	f, err := client.Files.Get(msg.Params.TextDocument.Uri)
	if err != nil {
		client.InternalError(msg.Id, err)
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
				client.Send(
					message.NewResponseMessage(
						msg.Id,
						map[string]interface{}{
							"uri": fmt.Sprintf("file:///%s", fname),
							"range": types.Range{
								Start: types.Position{Character: 0, Line: 0},
								End:   types.Position{Character: 0, Line: 0},
							},
						},
						nil,
					),
				)
				return
			} else {
				client.RequestFailed(msg.Id, fmt.Errorf("Could not find file %q", fname))
				return
			}
		}
	}
	client.RequestFailed(msg.Id, errors.New("No definition available"))
}

func TextDocumentHoverHandler(client client.Client, u message.UndefinedMessage) {
	msg := u.TextDocumentHover()

	f, err := client.Files.Get(msg.Params.TextDocument.Uri)
	if err != nil {
		client.InternalError(msg.Id, fmt.Errorf("Could not read file %s: %w", msg.Params.TextDocument.Uri, err))
		return
	}

	for _, t := range f.Tokens() {
		if t.Start.Line == msg.Params.Position.Line && t.Start.Character <= msg.Params.Position.Character && t.Stop.Character >= msg.Params.Position.Character {
			for k, category := range api.Sections {
				for _, v := range category {
					if t.Value == v.Name {
						client.Send(message.NewResponseMessage(msg.Id, map[string]any{
							"range": types.Range{
								Start: types.Position{Character: t.Start.Character, Line: t.Start.Line},
								End:   types.Position{Character: t.Stop.Character + 1, Line: t.Stop.Line},
							},
							"contents": map[string]any{
								"kind":  "markdown",
								"value": v.Markdown(),
							},
						}, nil))
						return
					}
				}
				if t.Value == k {
					names := []string{}
					for _, v := range category {
						names = append(names, v.Short())
					}
					client.Send(message.NewResponseMessage(msg.Id, map[string]any{
						"range": types.Range{
							Start: types.Position{Character: t.Start.Character, Line: t.Start.Line},
							End:   types.Position{Character: t.Stop.Character + 1, Line: t.Stop.Line},
						},
						"contents": map[string]any{
							"kind":  "markdown",
							"value": fmt.Sprintf("%s\n```lua\n%s\n```", k, strings.Join(names, "\n")),
						},
					}, nil))
					return
				}
			}
		}
	}
	client.RequestFailed(msg.Id, errors.New("No hover available"))
}

func TextDocumentCompletionHandler(client client.Client, u message.UndefinedMessage) {
	msg := u.Completion()
	items := []map[string]any{}

	category := api.Sections["/"]
	category = append(category, api.Sections["GlobalEvents"]...)
	if !strings.HasSuffix(msg.Params.TextDocument.Uri.Path(), "Global.-1.ttslua") {
		category = append(category, api.Sections["ObjectEvents"]...)
	}
	if msg.Params.Context.TriggerKind == types.TriggerCharacterTrigger {
		category = api.Sections["Object"]
	}

	for _, v := range category {
		kind := types.TextCompletion
		switch v.Kind {
		case apiParameterKindClass:
			kind = types.ClassCompletion
		case apiParameterKindConstant:
			kind = types.ConstantCompletion
		case apiParameterKindEvent:
			kind = types.EventCompletion
		case apiParameterKindFunction:
			kind = types.FunctionCompletion
		case apiParameterKindProperty:
			kind = types.PropertyCompletion
		}
		items = append(items, map[string]any{
			"label": v.Name,
			"kind":  kind,
			"documentation": map[string]any{
				"kind":  "markdown",
				"value": v.Markdown(),
			},
		})
	}

	client.Send(message.NewResponseMessage(msg.Id, map[string]any{
		"items": items,
	}, nil))
}

func TextDocumentSignatureHelpHandler(client client.Client, u message.UndefinedMessage) {
}

func TextDocumentCodeActionHandler(client client.Client, u message.UndefinedMessage) {
}
