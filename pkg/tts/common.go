package tts

import (
	"fmt"

	"github.com/lazzer64/ttsls/pkg/uri"
)

type ScriptState struct {
	Name   string `json:"name,omitempty"`
	Guid   string `json:"guid,omitempty"`
	Script string `json:"script,omitempty"`
}

func (ss ScriptState) URI() (uri.URI, error) {
	if ss.Guid == "" {
		return uri.Parse(ss.Name + ".ttslua")
	}
	return uri.Parse(fmt.Sprintf("%s.%s.ttslua", ss.Name, ss.Guid))
}
