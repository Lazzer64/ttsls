package tokens

import (
	"bytes"
	"fmt"
	"io"
)

type TokenType int

//go:generate stringer -type=TokenType
const (
	INVALID TokenType = iota

	AND
	BREAK
	DO
	ELSE
	ELSEIF
	END
	FALSE
	FOR
	FUNCTION
	GOTO
	IF
	IN
	LOCAL
	NIL
	NOT
	OR
	REPEAT
	RETURN
	THEN
	TRUE
	UNTIL
	WHILE
	INCLUDE

	PLUS
	MINUS
	MULT
	DIV
	MOD
	POW
	LEN

	BIT_AND
	BIT_XOR
	BIT_OR
	BIT_L_SHIFT
	BIT_R_SHIFT

	DIV_FLOOR
	EQUALITY
	INEQUALITY
	LESS_EQUAL
	GREAT_EQUAL
	LESS
	GREAT
	ASSIGN

	PAREN_L
	PAREN_R
	CURLY_L
	CURLY_R
	SQUARE_L
	SQUARE_R
	LABEL

	SEMI_COLON
	COLON
	COMMA
	DOT
	CONCAT
	VARIADIC

	COMMENT
	IDENTIFIER
	STRING
	NUMBER
)

type Token struct {
	Type        TokenType
	Value       string
	Start, Stop Position
}

func (t Token) String() string {
	return fmt.Sprintf(`[%s %q %v,%v]`, t.Type, t.Value, t.Start, t.Stop)
}

type Position struct {
	Line, Character int
}

var reserved = map[string]TokenType{
	"and":      AND,
	"break":    BREAK,
	"do":       DO,
	"else":     ELSE,
	"elseif":   ELSEIF,
	"end":      END,
	"false":    FALSE,
	"for":      FOR,
	"function": FUNCTION,
	"goto":     GOTO,
	"if":       IF,
	"in":       IN,
	"local":    LOCAL,
	"nil":      NIL,
	"not":      NOT,
	"or":       OR,
	"repeat":   REPEAT,
	"return":   RETURN,
	"then":     THEN,
	"true":     TRUE,
	"until":    UNTIL,
	"while":    WHILE,
}

var delim = map[rune]TokenType{
	'%': MOD,
	'&': BIT_AND,
	'(': PAREN_L,
	')': PAREN_R,
	'*': MULT,
	'+': PLUS,
	',': COMMA,
	';': SEMI_COLON,
	'[': SQUARE_L,
	']': SQUARE_R,
	'^': POW,
	'{': CURLY_L,
	'|': BIT_OR,
	'}': CURLY_R,
}

func Tokenize(b []byte) []Token {
	tokens := []Token{}
	b = bytes.ReplaceAll(b, []byte{'\r', '\n'}, []byte{'\n'})
	scan := newScanner(b)

	for {
		ch, err := scan.ScanRune()
		if err != nil {
			break
		}
		start := scan.Position()
		if ch == ' ' || ch == '\t' || ch == '\n' {
			// whitespace can be ignored
			continue
		}
		if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || ch == '_' {
			// valid identifiers and reserved words start with a letter or an underscore
			word, err := scan.ScanWord(ch)
			if err != nil {
				continue
			}
			var t *Token
			for w, typ := range reserved {
				if w == string(word) {
					t = &Token{typ, w, start, scan.Position()}
					break
				}
			}
			if t != nil {
				tokens = append(tokens, *t)
				continue
			}
			tokens = append(tokens, Token{IDENTIFIER, string(word), start, scan.Position()})
			continue
		}
		if ch >= '0' && ch <= '9' {
			// tokens starting with an integer must be a number
			num, err := scan.ScanNumber(ch)
			if err != nil {
				continue
			}
			tokens = append(tokens, Token{NUMBER, string(num), start, scan.Position()})
			continue
		}
		if ch == '"' || ch == '\'' {
			str, err := scan.ScanString(ch)
			if err != nil {
				continue
			}
			tokens = append(tokens, Token{STRING, string(str), start, scan.Position()})
			continue
		}
		if ch == '#' {
			ch, err := scan.ScanRune()
			if err != nil {
				continue
			}
			word, err := scan.ScanWord(ch)
			if err != nil {
				continue
			}
			if string(word) == "include" {
				inc, err := scan.ScanInclude()
				if err != nil {
					continue
				}
				tokens = append(tokens, Token{INCLUDE, string(inc), start, scan.Position()})
			} else {
				for i := 0; i < len(word); i++ {
					scan.UnscanRune()
				}
				tokens = append(tokens, Token{LEN, "#", start, scan.Position()})
			}
			continue
		}
		if ch == '-' {
			ch, _ := scan.ScanRune()
			if ch == '-' {
				comment, err := scan.ScanComment()
				if err != nil {
					continue
				}
				tokens = append(tokens, Token{COMMENT, string(comment), start, scan.Position()})
			} else {
				scan.UnscanRune()
				tokens = append(tokens, Token{MINUS, "-", start, scan.Position()})
				continue
			}
		}
		var t *Token
		for r, typ := range delim {
			if ch == r {
				t = &Token{typ, string(ch), start, scan.Position()}
				break
			}
		}
		if t != nil {
			tokens = append(tokens, *t)
			continue
		}
		if ch == '.' {
			if ch, _ := scan.ScanRune(); ch == '.' {
				if ch, _ := scan.ScanRune(); ch == '.' {
					tokens = append(tokens, Token{VARIADIC, "...", start, scan.Position()})
				} else {
					scan.UnscanRune()
					tokens = append(tokens, Token{CONCAT, "..", start, scan.Position()})
				}
			} else {
				scan.UnscanRune()
				tokens = append(tokens, Token{DOT, ".", start, scan.Position()})
			}
			continue
		}
		if ch == '/' {
			if ch, _ := scan.ScanRune(); ch == '/' {
				tokens = append(tokens, Token{DIV_FLOOR, "//", start, scan.Position()})
			} else {
				scan.UnscanRune()
				tokens = append(tokens, Token{DIV, "/", start, scan.Position()})
			}
			continue
		}
		if ch == ':' {
			if ch, _ := scan.ScanRune(); ch == ':' {
				tokens = append(tokens, Token{LABEL, "::", start, scan.Position()})
			} else {
				scan.UnscanRune()
				tokens = append(tokens, Token{COLON, ":", start, scan.Position()})
			}
			continue
		}
		if ch == '<' {
			ch, _ := scan.ScanRune()
			if ch == '<' {
				tokens = append(tokens, Token{BIT_L_SHIFT, "<<", start, scan.Position()})
			} else if ch == '=' {
				tokens = append(tokens, Token{LESS_EQUAL, "<=", start, scan.Position()})
			} else {
				scan.UnscanRune()
				tokens = append(tokens, Token{LESS, "<", start, scan.Position()})
			}
			continue
		}
		if ch == '=' {
			if ch, _ := scan.ScanRune(); ch == '=' {
				tokens = append(tokens, Token{EQUALITY, "==", start, scan.Position()})
			} else {
				scan.UnscanRune()
				tokens = append(tokens, Token{ASSIGN, "=", start, scan.Position()})
			}
			continue
		}
		if ch == '>' {
			ch, _ := scan.ScanRune()
			if ch == '>' {
				tokens = append(tokens, Token{BIT_R_SHIFT, ">>", start, scan.Position()})
			} else if ch == '=' {
				tokens = append(tokens, Token{GREAT_EQUAL, ">=", start, scan.Position()})
			} else {
				scan.UnscanRune()
				tokens = append(tokens, Token{GREAT, ">", start, scan.Position()})
			}
			continue
		}
		if ch == '~' {
			if ch, _ := scan.ScanRune(); ch == '=' {
				tokens = append(tokens, Token{INEQUALITY, "~=", start, scan.Position()})
			} else {
				scan.UnscanRune()
				tokens = append(tokens, Token{BIT_XOR, "~", start, scan.Position()})
			}
			continue
		}
	}

	return tokens
}

func TokenizeString(s string) []Token {
	return Tokenize([]byte(s))
}

type scanner struct {
	r []byte
	i int
}

func newScanner(b []byte) *scanner {
	return &scanner{b, -1}
}

func (s *scanner) ScanRune() (rune, error) {
	s.i++
	if s.i+1 > len(s.r) {
		return 0, io.EOF
	}
	return rune(s.r[s.i]), nil
}

func (s *scanner) UnscanRune() {
	if s.i > 0 {
		s.i--
	}
}

func (s *scanner) Position() Position {
	p := Position{0, 0}
	for i := 0; i < s.i; i++ {
		if s.r[i] == '\n' {
			p.Line++
			p.Character = 0
		} else {
			p.Character++
		}
	}
	return p
}

func (s *scanner) ScanWord(ch rune) (word []rune, err error) {
	for ('A' <= ch && ch <= 'Z') || ('a' <= ch && ch <= 'z') || ('0' <= ch && ch <= '9') || ch == '_' {
		word = append(word, ch)
		ch, err = s.ScanRune()
		if err != nil {
			return
		}
	}
	s.UnscanRune()
	return
}

func (s *scanner) ScanNumber(ch rune) (num []rune, err error) {
	// TODO: hex constants
	radix := false
	for ('0' <= ch && ch <= '9') || (ch == '.' && !radix) {
		if ch == '.' {
			radix = true
		}
		num = append(num, ch)
		ch, err = s.ScanRune()
		if err != nil {
			return
		}
	}
	s.UnscanRune()
	return
}

func (s *scanner) ScanString(ch rune) (str []rune, err error) {
	end := ch
	escape := false

	ch, err = s.ScanRune()
	for escape || ch != end {
		str = append(str, ch)
		ch, err = s.ScanRune()
		if err != nil {
			return
		}
		escape = false
		if ch == '\\' {
			escape = true
			ch, err = s.ScanRune()
			if err != nil {
				return
			}
		}
	}
	return
}

func (s *scanner) ScanComment() (comment []rune, err error) {
	// assumes the previous two characters are '--'
	// TODO: block comments
	var ch rune
	for {
		ch, err = s.ScanRune()
		if err != nil {
			return
		}
		if ch == '\n' {
			s.UnscanRune()
			break
		}
		comment = append(comment, ch)
	}
	return
}

func (s *scanner) ScanInclude() (fname []rune, err error) {
	// assumes the previous word was '#include'
	var ch rune
	for {
		ch, err = s.ScanRune()
		if err != nil {
			return
		}
		if ch != ' ' && ch != '\t' {
			break
		}
		if ch == '\n' {
			s.UnscanRune()
			return fname, fmt.Errorf("Unexpected end of line in #include statement")
		}
	}
	for ch != ' ' && ch != '\t' && ch != '\n' {
		fname = append(fname, ch)
		ch, err = s.ScanRune()
		if err != nil {
			return
		}
	}
	s.UnscanRune()
	return
}
