package handler

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/lazzer64/ttsls/pkg/lsp/client"
	"github.com/lazzer64/ttsls/pkg/lsp/protocol"
	"github.com/lazzer64/ttsls/pkg/tts"
)

func WorkspaceExecuteCommandHandler(c client.Client, u protocol.Message) {
	m := u.WorkspaceExecuteCommand()
	args := m.Params.Arguments.([]any)

	switch m.Params.Command {
	case "tts.exec":
		if len(args) != 1 {
			c.Send(protocol.NewErrorResponse(
				m.Id,
				protocol.ErrorCodesInvalidParams,
				fmt.Errorf("Wrong number of arguments. Expected 1 but got %d", len(args)),
			))
			break
		}

		if err := tts.Exec(fmt.Sprintf("%v", args[0])); err != nil {
			c.Send(protocol.NewErrorResponse(
				m.Id,
				protocol.ErrorCodesInternalError,
				err,
			))
		}

	case "tts.getScripts":
		if err := tts.GetScripts(); err != nil {
			c.Send(protocol.NewErrorResponse(
				m.Id,
				protocol.ErrorCodesInternalError,
				err,
			))
		}

	case "tts.saveAndPlay":
		states := []tts.ScriptState{}

		for _, sf := range c.Files.GetAllOpen() {
			parts := strings.Split(filepath.Base(sf.Path()), ".")
			if len(parts) >= 3 && parts[len(parts)-1] == "ttslua" && (parts[len(parts)-2] == "-1" || len(parts[len(parts)-2]) == 6) {
				script, err := client.Expand(c, sf)
				if err != nil {
					c.Log(protocol.MessageTypeError, err.Error())
					continue
				}
				states = append(states, tts.ScriptState{
					Name:   strings.Join(parts[0:len(parts)-2], "."),
					Guid:   parts[len(parts)-2],
					Script: script,
				})
			}
		}

		err := tts.SaveAndPlay(states...)

		if err != nil {
			c.Log(protocol.MessageTypeError, err.Error())
			break
		}

		c.Log(protocol.MessageTypeInfo, fmt.Sprintf(`Saved %d script`, len(states)))

	default:
		c.Send(protocol.NewErrorResponse(
			m.Id,
			protocol.ErrorCodesInvalidRequest,
			fmt.Errorf("Command not recognized"),
		))
	}
}
