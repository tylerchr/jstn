package jstn

import (
	"bytes"
	"strings"
	"testing"
)

func TestScanner(t *testing.T) {

	var buf bytes.Buffer

	s := NewScanner(strings.NewReader(WrittenCollectionSchema))
	for {
		tok, lit := s.Scan()
		if tok == EOF {
			break
		}
		if tok != WHITESPACE {
			t.Logf("%s %q\n", tok, lit)
			buf.WriteString(lit)
		}
	}

	t.Logf("Full: %q\n", buf.String())

}
