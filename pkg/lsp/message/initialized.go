package message

import (
	"encoding/json"
)

type InitializedMessage struct {
	RequestMessage
}

func (um UndefinedMessage) Initialized() (m InitializedMessage) {
	json.Unmarshal(um.raw, &m)
	return m
}
