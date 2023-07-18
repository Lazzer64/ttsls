package lsp

import (
	"github.com/lazzer64/ttsls/pkg/lsp/client"
	"github.com/lazzer64/ttsls/pkg/lsp/protocol"
	"github.com/lazzer64/ttsls/pkg/tts"
)

func (lsp *LSP) ttsHandler(c client.Client, msg tts.TTSMessage) {
	switch msg.MessageID {
	case tts.GAME_LOADED_RESPONSE:
		c.Log(protocol.MessageTypeInfo, "Game Loaded")

	case tts.NEW_OBJECT_RESPONSE:
		for _, ss := range msg.ScriptStates {
			uri, err := ss.URI()
			if err != nil {
				c.Log(protocol.MessageTypeError, err.Error())
				continue
			}
			c.Files.Write(uri, ss.Script)
			c.Send(
				protocol.NewWindowShowDocumentRequest(
					lsp.nextId(),
					protocol.ShowDocumentParams{
						Uri:       protocol.URI(uri),
						External:  false,
						TakeFocus: true,
					},
				),
			)
		}

	case tts.PRINT_RESPONSE:
		c.Log(protocol.MessageTypeInfo, msg.Message)

	case tts.ERROR_RESPONSE:
		c.Log(protocol.MessageTypeError, msg.Error)
	}
}
