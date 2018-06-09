package jstn

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func TestString(t *testing.T) {
	schema := Type{Kind: String, Optional: true}
	if expected, actual := "string?", schema.String(); expected != actual {
		t.Errorf("unexpected string: expected %q but got %q", expected, actual)
	}
}

func TestMarshalType(t *testing.T) {

	expected := []byte(`"string?"`)

	schema := Type{Kind: String, Optional: true}

	if data, err := json.Marshal(schema); err != nil {
		t.Errorf("unexpected marshaling error: %s", err)
	} else if !bytes.Equal(data, expected) {
		t.Errorf("unexpected marshal: expected %q but got %q", expected, data)
	}

}

func TestUnmarshalType(t *testing.T) {

	doc := []byte(`"string?"`)
	expected := Type{Kind: String, Optional: true}

	var s Type
	if err := json.Unmarshal(doc, &s); err != nil {
		t.Fatalf("unexpected unmarshaling error: %s", err)
	}

	if !reflect.DeepEqual(expected, s) {
		t.Errorf("unexpected unmarshal: expected %v but got %v", expected, s)
	}

}

func TestUnmarshalType_NotString(t *testing.T) {

	doc := []byte(`{"some":"object"}`)

	var s Type
	if err := json.Unmarshal(doc, &s); err == nil {
		t.Errorf("expected unmarshaling error but got type: %#v", s)
	}

}
