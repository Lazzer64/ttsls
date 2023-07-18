package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"
)

//go:embed templates/*.tmpl
var templatesDir embed.FS

func main() {
	if len(os.Args) != 2 {
		panic(fmt.Errorf("usage: %s <tag>", os.Args[0]))
	}
	resp, err := http.Get(fmt.Sprintf(
		"https://raw.githubusercontent.com/microsoft/vscode-languageserver-node/%s/protocol/metaModel.json",
		os.Args[1],
	))
	if err != nil {
		panic(err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var model MetaModel
	json.Unmarshal(b, &model)

	templates := template.New("templates").Funcs(template.FuncMap{
		"title":   strings.Title,
		"replace": strings.ReplaceAll,
		"comment": func(s string) string {
			if len(s) == 0 {
				return ""
			}
			return "// " + strings.Join(strings.Split(s, "\n"), "\n// ")
		},
		"indent": func(s string) string {
			if len(s) == 0 {
				return ""
			}
			return "\t" + strings.Join(strings.Split(s, "\n"), "\n\t")
		},
		"gotype": func(t Type) string {
			switch t.Name {
			case "boolean":
				return "bool"
			case "integer":
				return "int"
			case "decimal":
				return "float32"
			case "string":
				return "string"
			case "uinteger":
				return "uint32"
			case "":
				return "any"
			default:
				return t.Name
			}
		},
	})

	templates, err = templates.ParseFS(templatesDir, "templates/*.tmpl")
	if err != nil {
		panic(err)
	}
	for _, t := range templates.Templates() {
		name := fmt.Sprintf("%s_generated.go", strings.TrimSuffix(t.Name(), ".go.tmpl"))
		f, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}
		err = t.Execute(f, model)
		if err != nil {
			panic(err)
		}
	}
}

type MetaModel struct {
	MetaData      MetaData       `json:"metaData"`
	Requests      []Request      `json:"requests"`
	Notifications []Notification `json:"notifications"`
	Structures    []Structure    `json:"structures"`
	Enumerations  []Enumeration  `json:"enumerations"`
	TypeAliases   []TypeAlias    `json:"typeAliases"`
}

type MetaData struct {
	Version string `json:"version"`
}

type Request struct {
	Method              string              `json:"method"`
	Result              Result              `json:"result"`
	MessageDirection    MessageDirection    `json:"messageDirection"`
	Params              Type                `json:"params"`
	PartialResult       PartialResult       `json:"partialResult"`
	RegistrationOptions RegistrationOptions `json:"registrationOptions"`
	Documentation       string              `json:"documentation"`
}

type Result struct {
	Kind  ResultKind
	Items []Item
}

type ResultKind string

const (
	ResultKindOr ResultKind = "or"
)

type Item struct {
	Kind    string `json:"kind"`
	Name    string `json:"name"`
	Element *Item  `json:"element"`
}

type MessageDirection string

const (
	ClientToServer MessageDirection = "clientToServer"
	ServerToClient MessageDirection = "serverToClient"
	Both           MessageDirection = "both"
)

type PartialResult Result

type RegistrationOptions struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type Notification struct {
	Method              string           `json:"method"`
	MessageDirection    MessageDirection `json:"messageDirection"`
	Params              Type             `json:"params"`
	RegistrationOptions Type             `json:"registrationOptions"`
	Documentation       string           `json:"documentation"`
}

type Structure struct {
	Name          string     `json:"name"`
	Properties    []Property `json:"properties"`
	Extends       []Type     `json:"extends"`
	Mixins        []Type     `json:"mixins"`
	Documentation string     `json:"documentation"`
}

type Property struct {
	Name          string `json:"name"`
	Type          Type   `json:"type"`
	Optional      bool   `json:"optional"`
	Documentation string `json:"documentation"`
}

type Enumeration struct {
	Name                 string             `json:"name"`
	Type                 Type               `json:"type"`
	Values               []EnumerationValue `json:"values"`
	SupportsCustomValues bool               `json:"supportsCustomValues"`
	Documentation        string             `json:"documentation"`
}

type EnumerationValue struct {
	Name          string `json:"name"`
	Value         any    `json:"value"`
	Documentation string `json:"documentation"`
}

type Type struct {
	Kind  string `json:"kind"`
	Name  string `json:"name"`
	Items []Item `json:"items"`
}

type TypeAlias struct {
	Name          string `json:"name"`
	Type          Type   `json:"type"`
	Documentation string `json:"documentation"`
}
