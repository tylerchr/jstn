package jstn

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type token int

var eof = rune(0)

const (
	// Control
	ILLEGAL token = iota
	EOF
	WHITESPACE
	NEWLINE

	// Literals
	IDENT

	// Known identifiers
	STRING  // string
	NUMBER  // number
	BOOLEAN // boolean
	NULL    // null

	// Structural characters
	CURLYOPEN   // {
	CURLYCLOSE  // }
	SQUAREOPEN  // ]
	SQUARECLOSE // ]
	COLON       // :
	SEMICOLON   // ;
	QUESTION    // ?
)

func (t token) String() string {
	return tokens[t]
}

var tokens = map[token]string{
	ILLEGAL:     "ILLEGAL",
	EOF:         "EOF",
	WHITESPACE:  "WHITESPACE",
	NEWLINE:     "NEWLINE",
	IDENT:       "IDENT",
	CURLYOPEN:   "CURLYOPEN",
	CURLYCLOSE:  "CURLYCLOSE",
	SQUAREOPEN:  "SQUAREOPEN",
	SQUARECLOSE: "SQUARECLOSE",
	COLON:       "COLON",
	SEMICOLON:   "SEMICOLON",
	QUESTION:    "QUESTION",
	STRING:      "STRING",
	NUMBER:      "NUMBER",
	BOOLEAN:     "BOOLEAN",
	NULL:        "NULL",
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t'
}

func isNewline(ch rune) bool {
	return ch == '\r' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

type scanner struct {
	r *bufio.Reader
}

func newScanner(r io.Reader) *scanner {
	return &scanner{r: bufio.NewReader(r)}
}

func (s *scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *scanner) unread() { _ = s.r.UnreadRune() }

func (s *scanner) Scan() (tok token, lit string) {

	ch := s.read()

	switch {
	case isWhitespace(ch):
		s.unread()
		return s.scanWhitespace()
	case isNewline(ch):
		s.unread()
		return s.scanNewline()
	case isWhitespace(ch):
		s.unread()
		return s.scanWhitespace()
	case isLetter(ch):
		s.unread()
		return s.scanIdent()
	case ch == eof:
		return EOF, ""
	}

	chars := map[rune]token{
		'{': CURLYOPEN,
		'}': CURLYCLOSE,
		'[': SQUAREOPEN,
		']': SQUARECLOSE,
		':': COLON,
		';': SEMICOLON,
		'?': QUESTION,
	}

	if tok, ok := chars[ch]; ok {
		return tok, string(ch)
	}

	return ILLEGAL, string(ch)
}

func (s *scanner) scanRunes(t token, matcher func(rune) bool) (tok token, lit string) {

	var buf bytes.Buffer

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !matcher(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return t, buf.String()

}

func (s *scanner) scanWhitespace() (tok token, lit string) {
	return s.scanRunes(WHITESPACE, func(r rune) bool {
		return isWhitespace(r)
	})
}

func (s *scanner) scanNewline() (tok token, lit string) {
	return s.scanRunes(NEWLINE, func(r rune) bool {
		return isNewline(r)
	})
}

func (s *scanner) scanIdent() (tok token, lit string) {

	tok, lit = s.scanRunes(IDENT, func(r rune) bool {
		return isLetter(r) || isDigit(r) || r == '_'
	})

	switch strings.ToLower(lit) {
	case "string":
		return STRING, lit
	case "number":
		return NUMBER, lit
	case "boolean":
		return BOOLEAN, lit
	case "null":
		return NULL, lit
	}

	return tok, lit

}
