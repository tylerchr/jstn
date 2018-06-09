package jstn

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
)

// Valid indicates whether the JSON document in is considered valid with
// respect to the JSTN structure t.
func Valid(t Type, in json.RawMessage) bool {
	return valid(t, in, false)
}

// StrictlyValid indicates whether the JSON document in is considered valid with
// respect to the JSTN structure t by strict standards.
func StrictlyValid(t Type, in json.RawMessage) bool {
	return valid(t, in, true)
}

// valid validates the JSON value in with respect to t. Whether to check for
// validity or strict validity is indicated by strictMode.
func valid(t Type, in json.RawMessage, strictMode bool) bool {
	d := json.NewDecoder(bytes.NewReader(in))
	d.UseNumber()

	// assert that the next json object matches the type
	if ok := validValue(d, t, strictMode); !ok {
		return false
	}

	// assert that all data has been parsed
	if _, err := d.Token(); err != io.EOF {
		return false
	}

	return true
}

// validValue indicates whether the next JSON value in the Decoder has the structure
// described by t.
func validValue(d *json.Decoder, t Type, strictMode bool) bool {

	switch t.Kind {
	case String:
		if ok := validString(d, t); !ok {
			return false
		}
	case Number:
		if ok := validNumber(d, t); !ok {
			return false
		}
	case Boolean:
		if ok := validBoolean(d, t); !ok {
			return false
		}
	case Null:
		if ok := validNull(d, t); !ok {
			return false
		}
	case Array:
		if ok := validArray(d, t, strictMode); !ok {
			return false
		}
	case Object:
		if ok := validObject(d, t, strictMode); !ok {
			return false
		}
	case Any:
		if ok := validAny(d, t, strictMode); !ok {
			return false
		}
	default:
		return false
	}

	return true
}

// validString indicates whether the next token in the JSON is a valid string.
func validString(d *json.Decoder, t Type) bool {

	tok, err := d.Token()
	if err != nil {
		if t.Optional && err == io.EOF {
			// no token at all, but it's optional so that's okay
			return true
		}
		log.Printf("validation failed: string: got error: %s\n", err)
		return false
	}

	switch tok.(type) {
	case string:
		return true
	case nil:
		if t.Optional {
			return true
		} else {
			log.Println("validation failed: string: got nil for required value")
			return false
		}
	default:
		log.Printf("validation failed: string: unexpected type %T: %v\n", tok, tok)
		return false
	}

}

// validNumber indicates whether the next token in the JSON is a valid number.
func validNumber(d *json.Decoder, t Type) bool {

	tok, err := d.Token()
	if err != nil {
		if t.Optional && err == io.EOF {
			// no token at all, but it's optional so that's okay
			return true
		}
		log.Printf("validation failed: number: got error: %s\n", err)
		return false
	}

	switch tok.(type) {
	case json.Number:
		return true
	case nil:
		if t.Optional {
			return true
		} else {
			log.Println("validation failed: number: got nil for required value")
			return false
		}
	default:
		log.Printf("validation failed: number: unexpected type %T: %v\n", tok, tok)
		return false
	}

}

// validBoolean indicates whether the next token in the JSON is a valid boolean.
func validBoolean(d *json.Decoder, t Type) bool {

	tok, err := d.Token()
	if err != nil {
		if t.Optional && err == io.EOF {
			// no token at all, but it's optional so that's okay
			return true
		}
		log.Printf("validation failed: boolean: got error: %s\n", err)
		return false
	}

	switch tok.(type) {
	case bool:
		return true
	case nil:
		if t.Optional {
			return true
		} else {
			log.Println("validation failed: boolean: got nil for required value")
			return false
		}
	default:
		log.Printf("validation failed: boolean: unexpected type %T: %v\n", tok, tok)
		return false
	}

}

// validNull indicates whether the next token in the JSON is a valid null value.
func validNull(d *json.Decoder, t Type) bool {

	tok, err := d.Token()
	if err != nil {
		if t.Optional && err == io.EOF {
			// no token at all, but it's optional so that's okay
			return true
		}
		log.Printf("validation failed: null: got error: %s\n", err)
		return false
	}

	switch tok.(type) {
	case nil:
		return true
	default:
		log.Printf("validation failed: null: unexpected type %T: %v\n", tok, tok)
		return false
	}

}

// validArray indicates whether the next token in the JSON is a valid array.
func validArray(d *json.Decoder, t Type, strictMode bool) bool {

	tok, err := d.Token()
	if err != nil {
		if t.Optional && err == io.EOF {
			// no token at all, but it's optional so that's okay
			return true
		}
		log.Printf("validation failed: array: got error: %s\n", err)
		return false
	}

	if delim, ok := tok.(json.Delim); !ok || delim != json.Delim('[') {
		log.Printf("validation failed: array: unexpected delimiter: %s\n", delim)
		return false
	}

	// schedule the consumption of the ending ']'
	defer d.Token()

	for d.More() {
		if ok := validValue(d, *t.Items, strictMode); !ok {
			log.Println("validation failed: array: invalid sub-element")
			return false
		}
	}

	return true

}

// validObject indicates whether the next token in the JSON is a valid object.
func validObject(d *json.Decoder, t Type, strictMode bool) bool {

	tok, err := d.Token()
	if err != nil {
		if t.Optional && err == io.EOF {
			// no token at all, but it's optional so that's okay
			return true
		}
		log.Printf("validation failed: object: got error: %s\n", err)
		return false
	}

	if delim, ok := tok.(json.Delim); !ok || delim != json.Delim('{') {
		log.Printf("validation failed: object: unexpected delimiter: %s\n", delim)
		return false
	}

	// schedule the consumption of the ending '}'
	defer d.Token()

	// note which properties we need to find
	necessaryProps := make(map[string]bool)
	for k := range t.Properties {
		necessaryProps[k] = true
	}

	for d.More() {

		// parse out a key
		keyTok, err := d.Token()
		if err != nil {
			log.Printf("validation failed: object: failed to read key: %s\n", err)
			return false
		}

		// it better be a string
		keyTokStr, keyOk := keyTok.(string)
		if !keyOk {
			log.Printf("validation failed: object: non-string key: %s\n", keyTok)
			return false
		}

		// check whether this property was declared
		propType, ok := t.Properties[keyTokStr]
		if !ok {
			if strictMode {
				log.Printf("validation failed: object: contains undeclared property %q\n", keyTokStr)
				return false
			}

			// in non-strict mode, undeclared properties have an implicit type of any
			propType = &Type{Kind: Any}
		}

		// mark this property as visited
		delete(necessaryProps, keyTokStr)

		if ok := validValue(d, *propType, strictMode); !ok {
			log.Println("validation failed: object: invalid sub-element")
			return false
		}
	}

	// make sure that any not-located properties were optional
	for k := range necessaryProps {
		if opt := t.Properties[k].Optional; !opt {
			log.Printf("validation failed: object: missing required property: %q\n", k)
			return false
		}
	}

	return true

}

// validAny indicates whether the next token in the JSON is any valid JSON value.
func validAny(d *json.Decoder, t Type, strictMode bool) bool {

	// make sure we can fetch another JSON value from the decoder, but
	// we don't really care what it is.
	var v json.RawMessage
	if err := d.Decode(&v); err != nil {
		log.Printf("validation failed: any: got error: %s\n", err)
		return false
	}

	if strictMode {
		return false
	}

	return true

}
