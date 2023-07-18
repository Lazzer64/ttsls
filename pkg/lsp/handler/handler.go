package handler

import (
	"github.com/lazzer64/ttsls/pkg/lsp/protocol"
	"github.com/lazzer64/ttsls/pkg/lsp/client"
)

type Handler func(client.Client, protocol.Message)
