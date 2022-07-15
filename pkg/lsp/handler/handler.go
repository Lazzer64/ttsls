package handler

import (
	"github.com/lazzer64/ttsls/pkg/lsp/message"
	"github.com/lazzer64/ttsls/pkg/lsp/client"
)

type Handler func(client.Client, message.UndefinedMessage)
