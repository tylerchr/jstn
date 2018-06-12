package jstn

import (
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {

	cases := []struct {
		Schema string
		Parsed Type
	}{
		{
			Schema: `string`,
			Parsed: Type{Kind: String, Optional: false},
		},
		{
			Schema: `string?`,
			Parsed: Type{Kind: String, Optional: true},
		},
		{
			Schema: `number`,
			Parsed: Type{Kind: Number, Optional: false},
		},
		{
			Schema: `number?`,
			Parsed: Type{Kind: Number, Optional: true},
		},
		{
			Schema: `boolean`,
			Parsed: Type{Kind: Boolean, Optional: false},
		},
		{
			Schema: `boolean?`,
			Parsed: Type{Kind: Boolean, Optional: true},
		},
		{
			Schema: `null`,
			Parsed: Type{Kind: Null, Optional: false},
		},
		{
			Schema: `null?`,
			Parsed: Type{Kind: Null, Optional: true},
		},
		{
			Schema: `[number]?`,
			Parsed: Type{Kind: Array, Optional: true, Items: &Type{
				Kind: Number,
			}},
		},
		{
			Schema: `[]?`,
			Parsed: Type{Kind: Array, Optional: true, Items: nil},
		},
		{
			Schema: `[number?]`,
			Parsed: Type{Kind: Array, Items: &Type{
				Kind:     Number,
				Optional: true,
			}},
		},
		{
			Schema: `{key: string}`,
			Parsed: Type{Kind: Object, Properties: map[string]*Type{
				"key": &Type{Kind: String},
			}},
		},
		{
			Schema: `{}`,
			Parsed: Type{Kind: Object, Properties: map[string]*Type{}},
		},
		{
			Schema: `{name:string;age:number?}`,
			Parsed: Type{Kind: Object, Properties: map[string]*Type{
				"name": &Type{Kind: String},
				"age":  &Type{Kind: Number, Optional: true},
			}},
		},
		{
			Schema: `{
	name:string
	age:number?
}`,
			Parsed: Type{Kind: Object, Properties: map[string]*Type{
				"name": &Type{Kind: String},
				"age":  &Type{Kind: Number, Optional: true},
			}},
		},
		{
			Schema: `{
	name:string;
	age:number?;
}`,
			Parsed: Type{Kind: Object, Properties: map[string]*Type{
				"name": &Type{Kind: String},
				"age":  &Type{Kind: Number, Optional: true},
			}},
		},
		{
			Schema: `{author:string;works:[{
     title:string
     year:     number?;
     classic:boolean;}]}`,
			Parsed: Type{Kind: Object, Properties: map[string]*Type{
				"author": &Type{Kind: String},
				"works": &Type{Kind: Array, Items: &Type{Kind: Object, Properties: map[string]*Type{
					"title":   &Type{Kind: String},
					"year":    &Type{Kind: Number, Optional: true},
					"classic": &Type{Kind: Boolean},
				}}},
			}},
		},
	}

	for i, c := range cases {

		typedef, err := Parse(c.Schema)
		if err != nil {
			t.Fatalf("[case %d] failed to parse: %s\n", i, err)
		}

		if !reflect.DeepEqual(typedef, c.Parsed) {
			t.Errorf("[case %d] unexpected parse results: expected:\n%# v\nbut got\n%# v\n", i, c.Parsed, typedef)
		}

	}

}
