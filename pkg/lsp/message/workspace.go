package message

import (
	"encoding/json"
)

type WorkspaceExecuteCommandMessage struct {
	RequestMessage
	Params struct {
		Command   string   `json:"command"`
		Arguments []string `json:"arguments"`
	} `json:"params"`
}

func (um UndefinedMessage) WorkspaceExecuteCommand() (m WorkspaceExecuteCommandMessage) {
	json.Unmarshal(um.raw, &m)
	return m
}
