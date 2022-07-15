package client_test

import (
	"embed"
	"io"
	"testing"

	"github.com/lazzer64/ttsls/pkg/lsp/client"
)

//go:embed testdata
var testdata embed.FS

type testcase struct {
	name                string
	expanded, collapsed client.SourceFile
}

func Cases(t *testing.T, name string) (cases []testcase) {
	t.Helper()
	entries, err := testdata.ReadDir(name)
	if err != nil {
		t.Fatal(err)
	}
	for _, de := range entries {
		base := name + "/" + de.Name()
		expanded, err := client.ReadSourceFile(base + "/expanded.ttslua")
		if err != nil {
			t.Fatal(err)
		}
		collapsed, err := client.ReadSourceFile(base + "/collapsed.ttslua")
		if err != nil {
			t.Fatal(err)
		}
		cases = append(cases, testcase{
			name:      de.Name(),
			expanded:  expanded,
			collapsed: collapsed,
		})
	}
	return
}

func TestExpand(t *testing.T) {
	for _, tt := range Cases(t, "testdata/cases") {
		t.Run(tt.name, func(t *testing.T) {
			c := client.New(io.Discard)
			got, err := client.Expand(c, tt.collapsed)
			if err != nil {
				t.Error(err)
				return
			}
			want := tt.expanded.Content()
			if got != want {
				t.Errorf("Expand() = %v, want %v", got, want)
			}
		})
	}
}

func TestCollapse(t *testing.T) {
	for _, tt := range Cases(t, "testdata/cases") {
		t.Run(tt.name, func(t *testing.T) {
			c := client.New(io.Discard)
			got, err := client.Collapse(c, tt.expanded)
			if err != nil {
				t.Error(err)
				return
			}
			want := tt.collapsed.Content()
			if got != want {
				t.Errorf("Collapse() = %v, want %v", got, want)
			}
		})
	}
}
