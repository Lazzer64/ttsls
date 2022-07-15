package handler

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

//go:embed api.json
var data []byte

type apis struct {
	Sections apiSections
}

type apiSections map[string][]apiItem

type apiItem struct {
	Name        string           `json:"name"`
	Kind        apiParameterKind `json:"kind"`
	Type        string           `json:"type"`
	Description string           `json:"description"`
	Url         string           `json:"url"`
	Parameters  []apiParameter   `json:"parameters,omitempty"`
}

type apiParameter struct {
	Name       string         `json:"name"`
	Type       string         `json:"type"`
	TableItems []apiParameter `json:"table_items,omitempty"`
	Parameters []apiParameter `json:"parameters"`
}

func (p apiParameter) Markdown() string {
	if len(p.TableItems) > 0 {
		lines := []string{}
		for _, ti := range p.TableItems {
			lines = append(lines, fmt.Sprintf("**%s.%s**: %s", p.Name, ti.Name, ti.Type))
		}
		return strings.Join(lines, "\n")
	}
	return fmt.Sprintf("**%s**: %s", p.Name, p.Type)
}

type apiParameterKind string

const (
	apiParameterKindClass    apiParameterKind = "class"
	apiParameterKindConstant apiParameterKind = "constant"
	apiParameterKindEvent    apiParameterKind = "event"
	apiParameterKindFunction apiParameterKind = "function"
	apiParameterKindProperty apiParameterKind = "property"
)

func (itm apiItem) Markdown() string {
	if itm.Kind == apiParameterKindFunction || itm.Kind == apiParameterKindEvent {
		args := []string{}
		argDetails := []string{}
		for _, p := range itm.Parameters {
			args = append(args, p.Name)
			argDetails = append(argDetails, p.Markdown())
		}
		return strings.Join([]string{
			"```lua",
			fmt.Sprintf("function %s(%s) -- %s", itm.Name, strings.Join(args, ", "), itm.Type),
			"```",
			itm.Description,
			"",
			strings.Join(argDetails, "\n"),
			"",
			fmt.Sprintf("[docs](%s)", itm.Url),
		}, "\n")
	}

	return fmt.Sprintf("%s\n[docs](%s)", itm.Description, itm.Url)
}

func (itm apiItem) Short() string {
	if itm.Kind == apiParameterKindFunction || itm.Kind == apiParameterKindEvent {
		paramNames := []string{}
		for _, p := range itm.Parameters {
			paramNames = append(paramNames, p.Name)
		}
		return fmt.Sprintf("function %s(%s)", itm.Name, strings.Join(paramNames, ", "))
	}
	return itm.Name
}

var api apis

func init() {
	if err := json.Unmarshal(data, &api); err != nil {
		log.Fatal(err)
	}
}
