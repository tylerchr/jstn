package jstn

import (
	"encoding/json"
	"testing"
)

func TestValid(t *testing.T) {

	cases := []struct {
		Type     Type
		TestData json.RawMessage
		Valid    bool
	}{
		// A JSON string validates as a required string.
		{
			Type:     Type{Kind: KindString, Optional: false},
			TestData: json.RawMessage(`"a string"`),
			Valid:    true,
		},

		// An empty document validates as an optional string.
		{
			Type:     Type{Kind: KindString, Optional: true},
			TestData: json.RawMessage(``),
			Valid:    true,
		},

		// A JSON null validates as an empty string.
		{
			Type:     Type{Kind: KindString, Optional: true},
			TestData: json.RawMessage(`null`),
			Valid:    true,
		},

		// A JSON boolean does not validate as a required string.
		{
			Type:     Type{Kind: KindString, Optional: false},
			TestData: json.RawMessage(`true`),
			Valid:    false,
		},

		{
			Type:     Type{Kind: KindNumber, Optional: false},
			TestData: json.RawMessage(`1`),
			Valid:    true,
		},

		{
			Type:     Type{Kind: KindNumber, Optional: false},
			TestData: json.RawMessage(`-1`),
			Valid:    true,
		},

		{
			Type:     Type{Kind: KindNumber, Optional: false},
			TestData: json.RawMessage(`3.1415`),
			Valid:    true,
		},

		{
			Type:     Type{Kind: KindBoolean, Optional: false},
			TestData: json.RawMessage(`true`),
			Valid:    true,
		},

		{
			Type:     Type{Kind: KindBoolean, Optional: false},
			TestData: json.RawMessage(`false`),
			Valid:    true,
		},

		{
			Type:     Type{Kind: KindBoolean, Optional: true},
			TestData: json.RawMessage(``),
			Valid:    true,
		},

		{
			Type:     Type{Kind: KindBoolean, Optional: false},
			TestData: json.RawMessage(`3.1415`),
			Valid:    false,
		},

		{
			Type:     Type{Kind: KindNull, Optional: false},
			TestData: json.RawMessage(`null`),
			Valid:    true,
		},

		{
			Type:     Type{Kind: KindNull, Optional: true},
			TestData: json.RawMessage(`null`),
			Valid:    true,
		},

		{
			Type:     Type{Kind: KindNull, Optional: true},
			TestData: json.RawMessage(`"null"`),
			Valid:    false,
		},

		{
			Type: Type{Kind: KindArray, Optional: false, Items: &Type{
				Kind:     KindNumber,
				Optional: false,
			}},
			TestData: json.RawMessage(`[1, 2, 3, null]`),
			Valid:    false,
		},

		{
			Type: Type{Kind: KindArray, Optional: false, Items: &Type{
				Kind:     KindNumber,
				Optional: true,
			}},
			TestData: json.RawMessage(`[1, 2, 3, null]`),
			Valid:    true,
		},

		{
			Type: Type{Kind: KindObject, Optional: false, Properties: map[string]*Type{
				"foo": &Type{Kind: KindNull},
				"bar": &Type{Kind: KindNumber},
			}},
			TestData: json.RawMessage(`{"foo":null,"bar":123}`),
			Valid:    true,
		},

		{
			Type: Type{Kind: KindObject, Optional: false, Properties: map[string]*Type{
				"foo": &Type{Kind: KindNull},
				"bar": &Type{Kind: KindNumber},
			}},
			TestData: json.RawMessage(`{"foo":null,"bar":123,"baz":"hack"}`),
			Valid:    false,
		},

		{
			Type: Type{Kind: KindObject, Optional: false, Properties: map[string]*Type{
				"foo": &Type{Kind: KindNull},
				"bar": &Type{Kind: KindNumber},
			}},
			TestData: json.RawMessage(`{"foo":null}`),
			Valid:    false,
		},

		{
			Type: Type{Kind: KindObject, Optional: false, Properties: map[string]*Type{
				"foo": &Type{Kind: KindNull},
				"bar": &Type{Kind: KindNumber},
				"baz": &Type{Kind: KindString},
			}},
			TestData: json.RawMessage(`{"foo":null,"bar":123,"baz":"hack"}`),
			Valid:    true,
		},

		{
			Type: Type{Kind: KindObject, Optional: false, Properties: map[string]*Type{
				"foo": &Type{Kind: KindObject, Properties: map[string]*Type{
					"length": &Type{Kind: KindNumber},
				}},
				"bar": &Type{Kind: KindNumber},
			}},
			TestData: json.RawMessage(`{"foo":{"length":3},"bar":123}`),
			Valid:    true,
		},
	}

	for i, c := range cases {
		if ok := Valid(c.Type, c.TestData); ok != c.Valid {
			t.Errorf("[case %d] unexpected result for Valid: expected %t but got %t\n", i, c.Valid, ok)
		}
	}

}

func TestValidAPI(t *testing.T) {

	schema := MustParse(`{
	renderingOptions: {
		orientation: string?
	}
	inputs: [{
		inputId: string
		type: string
		value: number?
	}]
}
`)

	ok := Valid(schema, json.RawMessage(`{
	"renderingOptions": {},
	"inputs": [
		{
			"inputId": "some string",
			"type": "some type"
		}
	]
}`))

	t.Logf("Valid: %t\n", ok)

}

func BenchmarkValid(b *testing.B) {

	schema := MustParse(`{
	renderingOptions: {
		orientation: string?
	}
	inputs: [{
		inputId: string
		type: string
		value: number?
	}]
}
`)

	for i := 0; i < b.N; i++ {

		_ = Valid(schema, json.RawMessage(`{
		"renderingOptions": {},
		"inputs": [
			{
				"inputId": "some string",
				"type": "some type"
			}
		]
	}`))

	}

}
