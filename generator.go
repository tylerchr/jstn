package jstn

import (
	"bytes"
	"io"
	"sort"
	"strings"
)

// indentationString defines the whitespace string to use for indentation when
// generating using the pretty format.
const indentationString = "  "

// Generate formats t into a JSTN type declaration using the concise format
// defined by the JSTN specification.
func Generate(t Type) ([]byte, error) {
	return (generator{Pretty: false}).generate(t, 0), nil
}

// GeneratePretty formats t into a JSTN type declaration using the pretty
// format defined by the JSTN specification.
func GeneratePretty(t Type) ([]byte, error) {
	return (generator{Pretty: true, Indentation: indentationString}).generate(t, 0), nil
}

type generator struct {
	Pretty      bool   // Whether to render in pretty mode.
	Indentation string // When in pretty mode, the indentation character to use.
}

// generate formats t into a JSTN document. Because it is a recursive function,
// depth tracks the object hierarchy depth for use when pretty-printing.
func (g generator) generate(t Type, depth int) []byte {

	var buf bytes.Buffer

	switch t.Kind {
	case String:
		io.WriteString(&buf, "string") // token: string

	case Number:
		io.WriteString(&buf, "number") // token: number

	case Boolean:
		io.WriteString(&buf, "boolean") // token: boolean

	case Null:
		io.WriteString(&buf, "null") // token: null

	case Object:
		io.WriteString(&buf, "{") // token: begin-object

		// writePretty adds whitespace, but only if the generator is in pretty mode.
		writePretty := func(s string) {
			if g.Pretty && len(t.Properties) > 0 {
				io.WriteString(&buf, s)
			}
		}

		// Sort property names for determinism.
		var propertyNames []string
		for k := range t.Properties {
			propertyNames = append(propertyNames, k)
		}
		sort.Strings(propertyNames)

		writePretty("\n")
		for i, k := range propertyNames {

			// In pretty mode, indent the property declaration line.
			writePretty(strings.Repeat(g.Indentation, depth+1))

			// token: name
			io.WriteString(&buf, k)

			// token: name-separator
			io.WriteString(&buf, ":")
			writePretty(" ")

			// token: member
			buf.Write(g.generate(*t.Properties[k], depth+1))

			// token: delimiter
			writePretty("\n")
			if !g.Pretty && i < len(propertyNames)-1 {
				io.WriteString(&buf, ";")
			}
		}

		writePretty(strings.Repeat(g.Indentation, depth))
		io.WriteString(&buf, "}") // token: end-object

	case Array:
		io.WriteString(&buf, "[") // token: begin-array

		// token: type-declaration
		if t.Items != nil {
			buf.Write(g.generate(*t.Items, depth))
		}

		io.WriteString(&buf, "]") // token: end-array
	}

	if t.Optional {
		io.WriteString(&buf, "?") // token: value-optional
	}

	return buf.Bytes()

}
