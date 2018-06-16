package jstn

import "testing"

func TestGenerator(t *testing.T) {

	cases := []struct {
		Type   Type
		String string
		Pretty bool
	}{
		{
			Type:   Type{Kind: String, Optional: false},
			String: "string",
		},
		{
			Type:   Type{Kind: String, Optional: true},
			String: "string?",
		},
		{
			Type:   Type{Kind: Number, Optional: false},
			String: "number",
		},
		{
			Type:   Type{Kind: Number, Optional: true},
			String: "number?",
		},
		{
			Type:   Type{Kind: Boolean, Optional: false},
			String: "boolean",
		},
		{
			Type:   Type{Kind: Boolean, Optional: true},
			String: "boolean?",
		},
		{
			Type:   Type{Kind: Null, Optional: false},
			String: "null",
		},
		{
			Type:   Type{Kind: Null, Optional: true},
			String: "null?",
		},

		//
		// ARRAYS
		//

		{
			// This case is actually not permitted by the current spec, but is
			// proposed in https://github.com/tylerchr/jstn/issues/5.
			Type:   Type{Kind: Array, Optional: false, Items: nil},
			String: "[]",
		},
		{
			Type: Type{Kind: Array, Optional: false, Items: &Type{
				Kind: String,
			}},
			String: "[string]",
		},
		{
			// An array may be an optional type.
			Type: Type{Kind: Array, Optional: true, Items: &Type{
				Kind: String,
			}},
			String: "[string]?",
		},
		{
			// An array may contain an optional type.
			Type: Type{Kind: Array, Optional: false, Items: &Type{
				Kind: String, Optional: true,
			}},
			String: "[string?]",
		},

		//
		// OBJECTS
		//

		{
			Type:   Type{Kind: Object, Optional: false, Properties: nil},
			String: "{}",
		},
		{
			Type:   Type{Kind: Object, Optional: false, Properties: nil},
			String: "{}",
			Pretty: true,
		},
		{
			Type: Type{Kind: Object, Optional: false, Properties: map[string]*Type{
				"firstName": &Type{Kind: String, Optional: false},
				"age":       &Type{Kind: Number, Optional: true},
			}},
			String: "{age:number?;firstName:string}",
		},
		{
			Type: Type{Kind: Object, Optional: false, Properties: map[string]*Type{
				"firstName": &Type{Kind: String, Optional: false},
				"age":       &Type{Kind: Number, Optional: true},
			}},
			String: `{
  age: number?
  firstName: string
}`,
			Pretty: true,
		},
		{
			Type: Type{Kind: Object, Optional: true, Properties: map[string]*Type{
				"firstName": &Type{Kind: String, Optional: false},
				"age":       &Type{Kind: Number, Optional: true},
				"residences": &Type{Kind: Array, Optional: false, Items: &Type{
					Kind: Object, Optional: false, Properties: map[string]*Type{
						"city":    &Type{Kind: String, Optional: false},
						"country": &Type{Kind: String, Optional: true},
					},
				}},
			}},
			String: `{
  age: number?
  firstName: string
  residences: [{
    city: string
    country: string?
  }]
}?`,
			Pretty: true,
		},
	}

	for i, c := range cases {

		var out []byte
		var err error

		if c.Pretty {
			out, err = GeneratePretty(c.Type)
		} else {
			out, err = Generate(c.Type)
		}

		if err != nil {
			t.Errorf("[case %d] unexpected generator error: %s", i, err)
			continue
		}

		if string(out) != c.String {
			t.Errorf("[case %d] unexpected production:", i)
			t.Errorf(".             expected: %q", c.String)
			t.Errorf(".             got     : %q", out)
		}

	}

}
