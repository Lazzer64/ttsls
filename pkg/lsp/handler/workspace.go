package handler

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/lazzer64/ttsls/pkg/lsp/client"
	"github.com/lazzer64/ttsls/pkg/lsp/message"
	"github.com/lazzer64/ttsls/pkg/lsp/message/types"
	"github.com/lazzer64/ttsls/pkg/tts"
)

func WorkspaceExecuteCommandHandler(c client.Client, u message.UndefinedMessage) {
	msg := u.WorkspaceExecuteCommand()
	switch msg.Params.Command {
	case types.TTSExecCommand:
		if len(msg.Params.Arguments) != 1 {
			c.InvalidParams(
				msg.Id,
				errors.New(fmt.Sprintf("Wrong number of arguments. Expected 1 but got %d", len(msg.Params.Arguments))),
			)
			break
		}

		if err := tts.Exec(msg.Params.Arguments[0]); err != nil {
			c.InternalError(msg.Id, err)
		}

	case types.TTSGetScripts:
		if err := tts.GetScripts(); err != nil {
			c.InternalError(msg.Id, err)
		}

	case types.TTSSaveAndPlayCommand:
		states := []tts.ScriptState{}

		for _, sf := range c.Files.GetAllOpen() {
			parts := strings.Split(filepath.Base(sf.Path()), ".")
			if len(parts) >= 3 && parts[len(parts)-1] == "ttslua" && (parts[len(parts)-2] == "-1" || len(parts[len(parts)-2]) == 6) {
				script, err := client.Expand(c, sf)
				if err != nil {
					c.LogError(err.Error())
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
			c.InternalError(msg.Id, err)
			break
		}

		c.LogInfo(fmt.Sprintf(`Saved %d script`, len(states)))

	default:
		c.InvalidRequest(msg.Id, errors.New("Command not recongnized"))
	}
}
