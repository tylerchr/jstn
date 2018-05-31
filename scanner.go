package jstn

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type Token int

var eof = rune(0)

const (
	// Control
	ILLEGAL Token = iota
	EOF
	WHITESPACE
	NEWLINE

	// Literals
	IDENT

	// Misc characters
	CURLYOPEN   // {
	CURLYCLOSE  // }
	SQUAREOPEN  // ]
	SQUARECLOSE // ]
	COLON       // :
	SEMICOLON   // ;
	QUESTION    // ?

	// Known types
	STRING  // string
	NUMBER  // number
	BOOLEAN // boolean
	NULL    // null
)

func (t Token) String() string {
	return Tokens[t]
}

var Tokens = map[Token]string{
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

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() { _ = s.r.UnreadRune() }

func (s *Scanner) Scan() (tok Token, lit string) {

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

	chars := map[rune]Token{
		'{': CURLYOPEN,
		'}': CURLYCLOSE,
		'[': SQUAREOPEN,
		']': SQUARECLOSE,
		':': COLON,
		';': SEMICOLON,
		'?': QUESTION,
	}

	if token, ok := chars[ch]; ok {
		return token, string(ch)
	}

	return ILLEGAL, string(ch)
}

func (s *Scanner) scanRunes(t Token, matcher func(rune) bool) (tok Token, lit string) {

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

func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	return s.scanRunes(WHITESPACE, func(r rune) bool {
		return isWhitespace(r)
	})
}

func (s *Scanner) scanNewline() (tok Token, lit string) {
	return s.scanRunes(NEWLINE, func(r rune) bool {
		return isNewline(r)
	})
}

func (s *Scanner) scanIdent() (tok Token, lit string) {

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
