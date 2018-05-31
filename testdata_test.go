package jstn

var WrittenCollectionSchema = `{
	author: {
		penName: string?
	}
	works: [{
		title: string
		language: string
		pageCount: number?
	}]
}
`

var WrittenCollectionType = Type{
	Kind: KindObject,
	Properties: map[string]*Type{
		"author": &Type{
			Kind: KindObject,
			Properties: map[string]*Type{
				"penName": &Type{Kind: KindString, Optional: true},
			},
		},
		"works": &Type{
			Kind: KindArray,
			Items: &Type{
				Kind: KindObject,
				Properties: map[string]*Type{
					"title":     &Type{Kind: KindString},
					"language":  &Type{Kind: KindString},
					"pageCount": &Type{Kind: KindNumber, Optional: true},
				},
			},
		},
	},
}
