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
	Kind: Object,
	Properties: map[string]*Type{
		"author": &Type{
			Kind: Object,
			Properties: map[string]*Type{
				"penName": &Type{Kind: String, Optional: true},
			},
		},
		"works": &Type{
			Kind: Array,
			Items: &Type{
				Kind: Object,
				Properties: map[string]*Type{
					"title":     &Type{Kind: String},
					"language":  &Type{Kind: String},
					"pageCount": &Type{Kind: Number, Optional: true},
				},
			},
		},
	},
}
