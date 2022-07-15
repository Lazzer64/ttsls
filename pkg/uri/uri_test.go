package uri

import (
	"reflect"
	"runtime"
	"testing"
)

type UriTest struct {
	path string
	uri  URI
}

var winuritest = []UriTest{
	{
		"C:\\Users\\root\\Documents\\Global.-1.ttslua",
		"file:///C%3A/Users/root/Documents/Global.-1.ttslua",
	},
}

func TestParse(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skipf("uri.Parse: tests only available on windows")
	}
	for _, tt := range winuritest {
		t.Run("", func(t *testing.T) {
			got, err := Parse(tt.path)
			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.uri) {
				t.Errorf("Parse() = %v, want %v", got, tt.uri)
			}
		})
	}
}

func TestURI_Path(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skipf("uri.Parse: tests only available on windows")
	}
	for _, tt := range winuritest {
		t.Run("", func(t *testing.T) {
			if got := tt.uri.Path(); got != tt.path {
				t.Errorf("URI.Path() = %v, want %v", got, tt.path)
			}
		})
	}
}
