// Package jstn implements a reference parser and validator for JSON Type Notation.
package jstn

import (
	"fmt"
	"log"
	"strings"
)

func MustParse(schema string) Type {
	t, err := Parse(schema)
	if err != nil {
		panic(err)
	}
	return t
}

func Parse(schema string) (Type, error) {
	r := strings.NewReader(schema)
	p := &parser{s: newScanner(r)}
	return p.Parse()
}

type parser struct {
	s   *scanner
	buf struct {
		tok token
		lit string
		n   int
	}
}

func (p *parser) scan() (tok token, lit string) {
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}
	p.buf.tok, p.buf.lit = p.s.Scan()
	tok, lit = p.buf.tok, p.buf.lit
	return
}

func (p *parser) unscan() { p.buf.n = 1 }

func (p *parser) scanIgnoreWhitespace(ignoreNewlines bool) (tok token, lit string) {
	tok, lit = p.scan()
	for tok == WHITESPACE || (ignoreNewlines && tok == NEWLINE) {
		tok, lit = p.scan()
	}
	return
}

func (p *parser) Parse() (Type, error) {
	return p.parseType()
}

func (p *parser) parseType() (Type, error) {
	t, err := p.parseTypeDecl()
	if err != nil {
		return t, err
	}

	// if the next token is '?' mark it as optional
	tok, _ := p.scan()
	if tok == QUESTION {
		t.Optional = true
	} else {
		p.unscan()
	}

	return t, nil
}

func (p *parser) parseTypeDecl() (Type, error) {

	tok, _ := p.scanIgnoreWhitespace(true)

	switch tok {
	case STRING:
		return Type{Kind: String}, nil
	case NUMBER:
		return Type{Kind: Number}, nil
	case BOOLEAN:
		return Type{Kind: Boolean}, nil
	case NULL:
		return Type{Kind: Null}, nil
	case SQUAREOPEN:
		p.unscan()
		return p.parseArray()
	case CURLYOPEN:
		p.unscan()
		return p.parseObject()
	default:
		return Type{}, fmt.Errorf("unexpected token %s", tok.String())
	}
}

func (p *parser) parseArray() (Type, error) {

	// parse the opening brace
	tok, _ := p.scanIgnoreWhitespace(true)
	if tok != SQUAREOPEN {
		return Type{}, fmt.Errorf("unexpected token %s", tok.String())
	}

	// peek at the next character
	tok, _ = p.scanIgnoreWhitespace(true)
	p.unscan()

	var childType *Type
	if tok != SQUARECLOSE {

		// parse the internal type
		t, err := p.parseType()
		if err != nil {
			return Type{}, err
		}

		childType = &t

	}

	// parse the closing brace
	tok, _ = p.scanIgnoreWhitespace(true)
	if tok != SQUARECLOSE {
		return Type{}, fmt.Errorf("unexpected token %s", tok.String())
	}

	return Type{Kind: Array, Items: childType}, nil
}

func (p *parser) parseObject() (Type, error) {

	var tok token
	var lit string

	// parse the opening brace
	tok, _ = p.scanIgnoreWhitespace(true)
	if tok != CURLYOPEN {
		return Type{}, fmt.Errorf("unexpected token %s", tok.String())
	}

	props := make(map[string]*Type)

	for {

		tok, lit = p.scanIgnoreWhitespace(true)

		// check for the curly brace
		if tok == CURLYCLOSE {
			p.unscan()
			break
		}

		// parse the property name
		if tok != IDENT {
			log.Println("failing here")
			return Type{}, fmt.Errorf("unexpected token %s for %s", tok.String(), IDENT.String())
		}

		// parse the colon
		if tok, _ := p.scanIgnoreWhitespace(true); tok != COLON {
			return Type{}, fmt.Errorf("unexpected token %s for %s", tok.String(), COLON.String())
		}

		// parse the object type
		t, err := p.parseType()
		if err != nil {
			return Type{}, err
		}

		props[lit] = &t // save this property type

		// Every property pair must finish with a delimiter token, which can be
		// either a CURLYCLOSE (indicating the end of the object), a SEMICOLON,
		// or a NEWLINE.
		if tok, lit = p.scanIgnoreWhitespace(false); tok == CURLYCLOSE {

			// The SEMICOLON acts as a delimiter because it marks
			// the end of the entire object. There won't be any
			// more properties. Unscan it and let the next iteration
			// of the object property loop handle it.
			p.unscan()

		} else if tok != SEMICOLON && tok != NEWLINE {

			// If we didn't find a SEMICOLON, the only other valid
			// token is a NEWLINE. But we didn't find that either, so
			// we've got an error here.
			return Type{}, fmt.Errorf("unexpected token %s (%q) for [%s or %s]", tok.String(), lit, SEMICOLON.String(), NEWLINE.String())

		}

	}

	// parse the closing brace
	tok, _ = p.scanIgnoreWhitespace(true)
	if tok != CURLYCLOSE {
		return Type{}, fmt.Errorf("unexpected token %s", tok.String())
	}

	return Type{Kind: Object, Properties: props}, nil
}
