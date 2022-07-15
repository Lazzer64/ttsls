package tokens

import (
	"reflect"
	"testing"
)

const helloWorld = `-- print "Hello, World" and exit.
print("Hello, \"World\"")`

const includes = `#include foo
#include <bar>
#include !/baz
#include <~/bang>
`

func TestTokenize(t *testing.T) {
	type args struct {
		s []byte
	}
	tests := []struct {
		name       string
		args       args
		wantTokens []Token
	}{
		{
			"hello_world",
			args{[]byte(helloWorld)},
			[]Token{
				{COMMENT, ` print "Hello, World" and exit.`, Position{0, 0}, Position{0, 32}},
				{IDENTIFIER, "print", Position{1, 0}, Position{1, 4}},
				{PAREN_L, "(", Position{1, 5}, Position{1, 5}},
				{STRING, `Hello, "World"`, Position{1, 6}, Position{1, 23}},
				{PAREN_R, ")", Position{1, 24}, Position{1, 24}},
			},
		},
		{
			"includes",
			args{[]byte(includes)},
			[]Token{
				{INCLUDE, "foo", Position{0, 0}, Position{0, 11}},
				{INCLUDE, "<bar>", Position{1, 0}, Position{1, 13}},
				{INCLUDE, "!/baz", Position{2, 0}, Position{2, 13}},
				{INCLUDE, "<~/bang>", Position{3, 0}, Position{3, 16}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTokens := Tokenize(tt.args.s); !reflect.DeepEqual(gotTokens, tt.wantTokens) {
				t.Errorf("\ngot:  %v\nwant: %v", gotTokens, tt.wantTokens)
			}
		})
	}
}
