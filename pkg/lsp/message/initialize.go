package message

import (
	"encoding/json"

	"github.com/lazzer64/ttsls/pkg/lsp/message/types"
)

type InitializeMessage struct {
	RequestMessage
	Params types.InitializeParams `json:"params"`
}

func (um UndefinedMessage) Initialize() (m InitializeMessage) {
	json.Unmarshal(um.raw, &m)
	return m
}
