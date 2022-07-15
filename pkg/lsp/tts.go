package lsp

import (
	"github.com/lazzer64/ttsls/pkg/lsp/client"
	"github.com/lazzer64/ttsls/pkg/lsp/message"
	"github.com/lazzer64/ttsls/pkg/tts"
)

func (lsp *LSP) ttsHandler(c client.Client, msg tts.TTSMessage) {
	switch msg.MessageID {
	case tts.GAME_LOADED_RESPONSE:
		c.LogInfo("Game Loaded")

	case tts.NEW_OBJECT_RESPONSE:
		for _, ss := range msg.ScriptStates {
			uri, err := ss.URI()
			if err != nil {
				c.LogError(err.Error())
				continue
			}
			c.Files.Write(uri, ss.Script)
			c.Send(
				message.NewRequestMessage(
					lsp.nextId(),
					"window/showDocument",
					map[string]interface{}{
						"uri":       uri,
						"takeFocus": true,
					},
				),
			)
		}

	case tts.PRINT_RESPONSE:
		c.LogInfo(msg.Message)

	case tts.ERROR_RESPONSE:
		c.LogError(msg.Error)
	}
}
