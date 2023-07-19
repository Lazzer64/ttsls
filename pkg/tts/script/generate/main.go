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

//go:embed templates
var templatesDir embed.FS

func main() {
	if len(os.Args) != 2 {
		panic(fmt.Errorf("usage: %s <tag>", os.Args[0]))
	}
	resp, err := http.Get(fmt.Sprintf(
		"https://raw.githubusercontent.com/Berserk-Games/atom-tabletopsimulator-lua/%s/lib/api.json",
		os.Args[1],
	))
	if err != nil {
		panic(err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	data := map[string]any{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		panic(err)
	}

	templates := template.New("templates").Funcs(template.FuncMap{
		"title":   strings.Title,
		"replace": strings.ReplaceAll,
		"indent": func(s string) string {
			if len(s) == 0 {
				return ""
			}
			return "\t" + strings.Join(strings.Split(s, "\n"), "\n\t")
		},
	})

	templates, err = templates.ParseFS(templatesDir, "templates/*.tmpl", "templates/helpers/*")
	if err != nil {
		panic(err)
	}

	for _, t := range templates.Templates() {
		if !strings.HasSuffix(t.Name(), ".go.tmpl") {
			continue
		}
		name := fmt.Sprintf("%s_generated.go", strings.TrimSuffix(t.Name(), ".go.tmpl"))
		f, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}
		err = t.Execute(f, data)
		if err != nil {
			panic(err)
		}
	}
}
