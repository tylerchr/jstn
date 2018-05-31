package jstn

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
