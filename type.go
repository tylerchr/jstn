package jstn

import (
	"encoding/json"
)

// A Kind represents a primitive JSON type.
type Kind int

const (
	String Kind = iota
	Number
	Boolean
	Null
	Object
	Array
)

type Type struct {
	Kind       Kind
	Optional   bool
	Properties map[string]*Type // Only for Objects
	Items      *Type            // Only for Arrays
}

func (t Type) String() string {
	out, _ := Generate(t)
	return string(out)
}

// MarshalJSON implements json.Unmarshaler by converting t to its concise JSTN
// text representation.
func (t Type) MarshalJSON() ([]byte, error) {
	out, _ := Generate(t)
	return json.Marshal(string(out))
}

// UnmarshalJSON implements json.Unmarshaler by parsing a JSTN text string.
func (t *Type) UnmarshalJSON(b []byte) error {

	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	tt, err := Parse(s)
	*t = tt
	return err
}
