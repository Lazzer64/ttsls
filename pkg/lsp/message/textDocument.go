package message

import (
	"encoding/json"

	"github.com/lazzer64/ttsls/pkg/lsp/message/types"
)

type TextDocumentDefinitionMessage struct {
	RequestMessage
	Params struct {
		Position     types.Position               `json:"position"`
		TextDocument types.TextDocumentIdentifier `json:"textDocument"`
	} `json:"params"`
}

func (um UndefinedMessage) TextDocumentDefinition() (m TextDocumentDefinitionMessage) {
	json.Unmarshal(um.raw, &m)
	return m
}

type TextDocumentHoverMessage struct {
	RequestMessage
	Params struct {
		Position     types.Position               `json:"position"`
		TextDocument types.TextDocumentIdentifier `json:"textDocument"`
	} `json:"params"`
}

func (um UndefinedMessage) TextDocumentHover() (m TextDocumentHoverMessage) {
	json.Unmarshal(um.raw, &m)
	return m
}

type TextDocumentDidOpenMessage struct {
	RequestMessage
	Params struct {
		TextDocument types.TextDocumentItem `json:"textDocument"`
	} `json:"params"`
}

func (um UndefinedMessage) TextDocumentDidOpenMessage() (m TextDocumentDidOpenMessage) {
	json.Unmarshal(um.raw, &m)
	return m
}

type TextDocumentDidCloseMessage struct {
	RequestMessage
	Params struct {
		TextDocument types.TextDocumentIdentifier `json:"textDocument"`
	} `json:"params"`
}

func (um UndefinedMessage) TextDocumentDidCloseMessage() (m TextDocumentDidCloseMessage) {
	json.Unmarshal(um.raw, &m)
	return m
}

type TextDocumentDidChangeMessage struct {
	RequestMessage
	Params struct {
		TextDocument   types.TextDocumentIdentifier           `json:"textDocument"`
		ContentChanges []types.TextDocumentContentChangeEvent `json:"contentChanges"`
	} `json:"params"`
}

func (um UndefinedMessage) TextDocumentDidChangeMessage() (m TextDocumentDidChangeMessage) {
	json.Unmarshal(um.raw, &m)
	return m
}

type CodeActionMessage struct {
	RequestMessage
	Params struct {
		TextDocument types.TextDocumentIdentifier `json:"textDocument"`
		Range        types.Range                  `json:"range"`
	} `json:"params"`
}

func (um UndefinedMessage) CodeAction() (m CodeActionMessage) {
	json.Unmarshal(um.raw, &m)
	return m
}

type CompletionMessage struct {
	RequestMessage
	Params struct {
		TextDocument types.TextDocumentIdentifier `json:"textDocument"`
		Position     types.Position               `json:"position"`
		Context      struct {
			TriggerKind types.CompletionTriggerKind `json:"triggerKind"`
		} `json:"context"`
	} `json:"params"`
}

func (um UndefinedMessage) Completion() (m CompletionMessage) {
	json.Unmarshal(um.raw, &m)
	return m
}

type SignatureHelpMessage struct {
	RequestMessage
	Params struct {
		TextDocument types.TextDocumentIdentifier `json:"textDocument"`
		Position     types.Position               `json:"position"`
		Context      struct {
			TriggerKind         types.SignatureHelpTriggerKind `json:"triggerKind"`
			TriggerCharacter    string                         `json:"triggerCharacter"`
			IsRetrigger         bool                           `json:"isRetrigger"`
			ActiveSignatureHelp types.SignatureHelp            `json:"activeSignatureHelp"`
		} `json:"context"`
	} `json:"params"`
}

func (um UndefinedMessage) SignatureHelp() (m SignatureHelpMessage) {
	json.Unmarshal(um.raw, &m)
	return m
}
