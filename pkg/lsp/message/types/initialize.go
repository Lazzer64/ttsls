package types

import (
	"encoding/json"
)

type PositionEncodingKind string

const (
	PositionEncodingKindUTF8  PositionEncodingKind = "utf-8"
	PositionEncodingKindUTF16 PositionEncodingKind = "utf-16"
	PositionEncodingKindUTF32 PositionEncodingKind = "utf-32"
)

type ServerCapabilities struct {
	CompletionProvider     CompletionOptions       `json:"completionProvider,omitempty"`
	DefinitionProvider     DefinitionOptions       `json:"definitionProvider,omitempty"`
	ExecuteCommandProvider ExecuteCommandOptions   `json:"executeCommandProvider,omitempty"`
	Experimental           json.RawMessage         `json:"experimental,omitempty"`
	HoverProvider          HoverOptions            `json:"hoverProvider,omitempty"`
	PositionEncoding       PositionEncodingKind    `json:"positionEncoding,omitempty"`
	SignatureHelpProvider  SignatureHelpOptions    `json:"signatureHelpProvider,omitempty"`
	TextDocumentSync       TextDocumentSyncOptions `json:"textDocumentSync,omitempty"`
}

type TextDocumentSyncKind int

const (
	TextDocumentSyncKindNone        = 0
	TextDocumentSyncKindFull        = 1
	TextDocumentSyncKindIncremental = 2
)

type TextDocumentSyncOptions struct {
	OpenClose bool                 `json:"openClose,omitempty"`
	Change    TextDocumentSyncKind `json:"change,omitempty"`
}

type DefinitionOptions struct {
	DocumentSelector []DocumentFilter `json:"documentSelector"`
}

type HoverOptions struct {
	DocumentSelector []DocumentFilter `json:"documentSelector"`
}

type DocumentFilter struct {
	Language string `json:"language,omitempty"`
	Scheme   string `json:"scheme,omitempty"`
	Pattern  string `json:"pattern,omitempty"`
}

type ExecuteCommandOptions struct {
	Commands []string `json:"commands"`
}

type CompletionOptions struct {
	TriggerCharacters   []string `json:"triggerCharacters,omitempty"`
	AllCommitCharacters []string `json:"allCommitCharacters,omitempty"`
	ResolveProvider     bool     `json:"resolveProvider,omitempty"`
	CompletionItem      struct {
		LabelDetailsSupport bool `json:"labelDetailsSupport,omitempty"`
	} `json:"completionItem,omitempty"`
}

type SignatureHelpOptions struct {
	TriggerCharacters   []string `json:"triggerCharacters,omitempty"`
	RetriggerCharacters   []string `json:"retriggerCharacters,omitempty"`
}
