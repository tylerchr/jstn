package jstn

import (
	"reflect"
	"strings"
	"testing"
)

func TestScanner(t *testing.T) {

	var tokens []token

	s := newScanner(strings.NewReader(WrittenCollectionSchema))
	for {
		tok, _ := s.Scan()
		if tok == EOF {
			break
		}
		tokens = append(tokens, tok)
	}

	expected := []token{
		CURLYOPEN,
		NEWLINE,
		WHITESPACE,
		IDENT,
		COLON,
		WHITESPACE,
		CURLYOPEN,
		NEWLINE,
		WHITESPACE,
		IDENT,
		COLON,
		WHITESPACE,
		STRING,
		QUESTION,
		NEWLINE,
		WHITESPACE,
		CURLYCLOSE,
		NEWLINE,
		WHITESPACE,
		IDENT,
		COLON,
		WHITESPACE,
		SQUAREOPEN,
		CURLYOPEN,
		NEWLINE,
		WHITESPACE,
		IDENT,
		COLON,
		WHITESPACE,
		STRING,
		NEWLINE,
		WHITESPACE,
		IDENT,
		COLON,
		WHITESPACE,
		STRING,
		NEWLINE,
		WHITESPACE,
		IDENT,
		COLON,
		WHITESPACE,
		NUMBER,
		QUESTION,
		NEWLINE,
		WHITESPACE,
		CURLYCLOSE,
		SQUARECLOSE,
		NEWLINE,
		CURLYCLOSE,
		NEWLINE,
	}

	if !reflect.DeepEqual(expected, tokens) {
		t.Errorf("unexpected tokens: expected %v but got %v\n", expected, tokens)
	}

}
