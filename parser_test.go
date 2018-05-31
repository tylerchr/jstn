package jstn

import (
	"reflect"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {

	cases := []struct {
		Schema string
		Parsed Type
	}{
		{
			Schema: `string`,
			Parsed: Type{Kind: KindString, Optional: false},
		},
		{
			Schema: `string?`,
			Parsed: Type{Kind: KindString, Optional: true},
		},
		{
			Schema: `number`,
			Parsed: Type{Kind: KindNumber, Optional: false},
		},
		{
			Schema: `number?`,
			Parsed: Type{Kind: KindNumber, Optional: true},
		},
		{
			Schema: `boolean`,
			Parsed: Type{Kind: KindBoolean, Optional: false},
		},
		{
			Schema: `boolean?`,
			Parsed: Type{Kind: KindBoolean, Optional: true},
		},
		{
			Schema: `null`,
			Parsed: Type{Kind: KindNull, Optional: false},
		},
		{
			Schema: `null?`,
			Parsed: Type{Kind: KindNull, Optional: true},
		},
		{
			Schema: `[number]?`,
			Parsed: Type{Kind: KindArray, Optional: true, Items: &Type{
				Kind: KindNumber,
			}},
		},
		{
			Schema: `[number?]`,
			Parsed: Type{Kind: KindArray, Items: &Type{
				Kind:     KindNumber,
				Optional: true,
			}},
		},
		{
			Schema: `{key: string}`,
			Parsed: Type{Kind: KindObject, Properties: map[string]*Type{
				"key": &Type{Kind: KindString},
			}},
		},
		{
			Schema: `{}`,
			Parsed: Type{Kind: KindObject, Properties: map[string]*Type{}},
		},
		{
			Schema: `{name:string;age:number?}`,
			Parsed: Type{Kind: KindObject, Properties: map[string]*Type{
				"name": &Type{Kind: KindString},
				"age":  &Type{Kind: KindNumber, Optional: true},
			}},
		},
		{
			Schema: `{
	name:string
	age:number?
}`,
			Parsed: Type{Kind: KindObject, Properties: map[string]*Type{
				"name": &Type{Kind: KindString},
				"age":  &Type{Kind: KindNumber, Optional: true},
			}},
		},
		{
			Schema: `{
	name:string;
	age:number?;
}`,
			Parsed: Type{Kind: KindObject, Properties: map[string]*Type{
				"name": &Type{Kind: KindString},
				"age":  &Type{Kind: KindNumber, Optional: true},
			}},
		},
		{
			Schema: `{author:string;works:[{
     title:string
     year:     number?;
     classic:boolean;}]}`,
			Parsed: Type{Kind: KindObject, Properties: map[string]*Type{
				"author": &Type{Kind: KindString},
				"works": &Type{Kind: KindArray, Items: &Type{Kind: KindObject, Properties: map[string]*Type{
					"title":   &Type{Kind: KindString},
					"year":    &Type{Kind: KindNumber, Optional: true},
					"classic": &Type{Kind: KindBoolean},
				}}},
			}},
		},
	}

	for i, c := range cases {

		typedef, err := NewParser(strings.NewReader(c.Schema)).Parse()
		if err != nil {
			t.Fatalf("[case %d] failed to parse: %s\n", i, err)
		}

		if !reflect.DeepEqual(typedef, c.Parsed) {
			t.Errorf("[case %d] unexpected parse results: expected:\n%# v\nbut got\n%# v\n", i, c.Parsed, typedef)
		}

	}

}
