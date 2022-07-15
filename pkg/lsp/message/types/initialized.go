package types

import (
	"encoding/json"

	"github.com/lazzer64/ttsls/pkg/uri"
)

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

type TraceValue string

const (
	TraceValueOff      = "off"
	TraceValueMessages = "messages"
	TraceValueVerbose  = "verbose"
)

type InitializeParams struct {
	Capabilities          ClientCapabilities    `json:"capabilities"`
	ClientInfo            ClientInfo            `json:"clientInfo,omitempty"`
	InitializationOptions InitializationOptions `json:"initializationOptions,omitempty"`
	Locale                string                `json:"locale,omitempty"`
	ProcessId             int                   `json:"processId,omitempty"`
	RootUri               uri.URI            `json:"rootUri"`
	Trace                 TraceValue            `json:"traceValue,omitempty"`
	WorkspaceFolders      []WorkspaceFolder     `json:"workspaceFolders,omitempty"`
}

type WorkspaceFolder struct {
	Uri  uri.URI `json:"uri"`
	Name string     `json:"name"`
}

type InitializationOptions json.RawMessage

type ClientCapabilities json.RawMessage
