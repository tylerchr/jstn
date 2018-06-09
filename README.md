# `jstn` [![GoDoc](https://godoc.org/github.com/tylerchr/jstn?status.svg)](https://godoc.org/github.com/tylerchr/jstn) [![Build Status](https://travis-ci.org/tylerchr/jstn.svg?branch=master)](https://travis-ci.org/tylerchr/jstn)

Package `jstn` implements a reference parser and validator for JSON Type Notation.

A JSTN document describes the structure of a JSON document in a typesafe way. It can be used to communicate expectations about JSON documents and is particularly useful for validating JSON.

## Overview

A JSTN type declaration describes the structure of a JSON document and its types by mirroring the structure of that JSON document and making assertions about types. An example JSTN type declaration looks like this:

```
{
	author: string
	works:[{
		title: string
		year: number?
		classic: boolean
	}]
}
```

The following JSON document would be considered valid with respect to this type declaration.

```json
{
	"author": "Johann Wolfgang von Goethe",
	"works": [
		{
			"title": "Prometheus",
			"year": 1773,
			"classic": false
		},
		{
			"title": "Die Leiden des jungen Werthers",
			"classic": true
		},
		{
			"title": "Faust",
			"year": 1775,
			"classic": true
		},
	]
}
```

For a formal specification of JSTN, see [spec.md](https://github.com/tylerchr/jstn/blob/master/SPEC.md).