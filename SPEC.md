# JSON Type Notation

JSON Type Notation (JSTN) is a lightweight, text-based, language-independent
data definition format designed to describe the type of JSON objects.

## Status of This Document

This document is an early draft and should not be considered either complete
or static.

## Introduction

JSON Type Notation (JSTN) is a text format for declaring type information
about JSON values, which are described in [RFC7159](https://tools.ietf.org/html/rfc7159). JSTN is designed
such that a JSTN text mirrors the structure of JSON texts that satisfy the
type, but is just enough unlike JSON to discourage the representation of
JSTN texts as actual JSON documents.

### Conventions Used in This Document

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC2119](https://tools.ietf.org/html/rfc2119).

The grammatical rules in this document are to be interpreted as
described in [RFC5234](https://tools.ietf.org/html/rfc5234).

## JSTN Grammar

A JSTN text is a sequence of tokens. The set of tokens includes seven
structural characters, identifiers, and five specific identifier literals.

```
   JSTN-text = ws type ws
```

These are the seven structural characters:

```
   begin-array     = ws %x5B ws  ; [ left square bracket

   begin-object    = ws %x7B ws  ; { left curly bracket

   end-array       = ws %x5D ws  ; ] right square bracket

   end-object      = ws %x7D ws  ; } right curly bracket

   name-separator  = ws %x3A ws  ; : colon

   value-separator = ws %x3B ws  ; ; semicolon

   value-optional  = ws %x3F ws  ; ? question mark
```

Insignificant whitespace is allowed before or after any of the seven
structural characters.

```
      ws  = horizontal-ws / nl

      horizontal-ws = *(
               %x20 /              ; Space
               %x09 )              ; Horizontal tab

      nl  = *(
               %x0A /              ; Line feed or New line
               %x0D )              ; Carriage return
```

## Types

A JSTN type MUST be an object, array, or one of the following five type
literals:

```
   string number boolean null any
```

The literal names MUST be lowercase. No other literal names are allowed.

Types may be marked as optional by suffixing the object, array, or type
literal with a question mark character.

```
   type-declaration = concrete-type [ value-optional ]

   concrete-type    = object / array / string / number / boolean / null

   string           = %x73.74.72.69.6e.67      ; string

   number           = %x6e.75.6d.62.65.72      ; number

   boolean          = %x62.6f.6f.6c.65.61.6e   ; boolean

   null             = %x6e.75.6c.6c            ; null

   any              = %x61.6e.79               ; any
```

## Objects

An object structure is represented as a pair of curly brackets surrounding zero
or more name/type pairs (or members). A single colon separates each name from its
value. A value is followed by a delimiter if it is not the last pair in the
object. For aesthetic reasons, the newline token may be used as a delimiter to
allow the semicolon to be omitted in multiline texts.

```
   object    = begin-object [ member *( delimiter member ) ] end-object

   member    = name name-separator type-declaration

   delimiter = value-separator / nl

   name      = (
   	             %x30-39 /     ; 0-9
   	             %x41-5A /     ; A-Z
   	             %x61-7A )     ; a-z

```

Like JSON texts, the behavor of applications that consume JSTN objects with
non-unique keys is unpredictable.

## Arrays

An array structure is represented as a pair of square brackets surrounding a
single type declaration. This single type declaration MUST be interpreted as
the type declaration for all elements contained within a validating JSON array.

This implies that while it is the case for JSON arrays described in RFC 7159
that

  "There is no requirement that the values in an array be of the same type."

such texts are not representable by JSTN, except through the usage of the
wildcard `any` type as the array's inner type declaration.

```
   array = begin-array type-declaration end-array
```

## The `any` type

JSTN allows for the specification of a property with a type key of `any`. This
designates the existence of the given property, without making any claims about that
property's actual type. That is, a property declared with the type indicator of `any` may
be used to represent any valid JSON value—be it an object, array, string, number, boolean,
or null.

Note that properties declared with the `any` type may also be designated as optional by
suffixing the type literal with the optional (`?`) token. This behaves consistently with
the usage of the optional token when appended to any other type. The value represented
by the `any?` type may appear in the JSON document as any valid JSON type, or may be
omitted. Contrast that with a value represented by the non-optional `any` type, which
may appear in the JSON document as either an object, array, string, number, boolean, or null
value, but must not be omitted.

## Parsers

A JSTN parser transforms a JSTN text into another representation. A JSTN
parser MUST accept all texts that conform to the JSTN grammar.

## Generators

A JSTN generator produces JSTN text. The resulting text MUST strictly conform
to the JSTN grammar.

A JSTN generator SHOULD provide mechanisms for generating the JSTN text in both
concise and pretty formats. The concise format omits all newlines (so object
pairs must be delineated by semicolons) but MAY include horizontal whitespace.
The pretty format includes newline characters (1) after each begin-object token
and (2) both prior to and after each end-object token. In the pretty format,
each line MUST be indented with an amount of whitespace corresponding to its
depth in the object hierarchy.

Other valid formatting variations exist, and a JSTN generator MAY additionally
implement support for other such variations.

## Validators

A JSTN validator accepts a JSTN type declaration and a JSON text, and indicates
whether the JSON text satisfies the JSTN type declaration. A JSON text is said
to be valid against a JSTN text if the following conditions all apply:

1. The JSON value is of the same type as the JSTN type declaration. This
   is applied recursively, such that a JSON object's properties must match
   the type of the same property at the same location in the JSTN type. A JSTN
   validator MUST fail validation if it encounters any values in the JSON
   document that do not match their declared JSTN type.

2. All non-optional types in the JSTN type declaration are present in the
   JSON document. A JSTN validator MUST fail validation upon detecting that
   a non-optional property is not present in the JSON document.

3. All optional types in the JSTN type declarations correspond either to
   (1) a value in the JSON document with a JSON type matching the JSTN type
   preceding the optional token for that type, (2) a JSON null value, or
   (3) in the case of object properties, that property's lack of presence.

A JSON document that does not satisfy these conditions with respect to a JSTN
text MUST NOT be considered valid with respect to that JSTN text.

### Strict Mode Validation

In addition to the standard validation conditions, A JSTN validator MAY choose to
offer an additional "strict mode" validator implementation. A strict JSTN validator
accepts a JSTN type declaration and a JSON text, enforces all of the validation
criteria as listed for a standard JSTN validator, and also validates the
JSON document against the following _additional_ criteria:

4. No object properties exist in the JSON document that are not declared in
   the JSTN type declaration text. A JSTN validator in strict mode MUST fail
   validation if properties are found in the JSON document that are not declared
   in the corresponding JSTN text.

5. The JSON document does not contain any properties that are represented by the
   JSTN text with an `any` or `any?` type. A JSTN validator in strict mode MUST
   fail upon encountering a JSON value that is declared in the JSTN text as having
   a type of `any` or `any?`.

A JSON document that does not satisfy these additional "strict mode" conditions with
respect to a JSTN text MAY be considered _invalid by "Strict Mode" JSTN standards_
with respect to that JSTN text. If the JSON document otherwise satisfies all of the
standard JSTN validation conditions, and only fails validation due to strict mode
conditions, usages of the validator SHOULD clarify that the document is considered
invalid only by _strict mode_ JSTN standards.

## Examples

This is a JSTN object represented in pretty format:

```
    {
        Image: {
            Width: number
            Height: number
            Title:  string
            License: string?
            Thumbnail: {
                Url:    string
                Height: number
                Width:  number
            }
            Animated: boolean?
            IDs: [number]
        }
    }
```

It describes the type of a JSON object whose Image member is an object whose
Thumbnail member is an object, whose IDs member is an array of numbers, and
whose Animated member is an optional boolean. The [first example in RFC 7159](https://tools.ietf.org/html/rfc7159#section-13)
is such an object.

Here is the same JSTN object represented in concise format:

```
   {Image:{Width:number;Height:number;Title:string;License:string?;Thumbnail:{Url:string;Format:string?;Height:number;Width:number};Animated:boolean?;IDs:[number]}}
```

This is a JSTN array represented in pretty format:

```
	[{
		precision: string
		Latitude: number
		Longitude: number
		Address: string
		City: string
		State: string
		Zip: string
		Country: string
		Planet: string?
	}]
```

The [second example in RFC 7159](https://tools.ietf.org/html/rfc7159#section-13) is considered valid with respect to this JSTN
type declaration.


Here is an example of a type declaration with more complexity:

```
   {
      userId: string
      firstName: string
      middleName: string?
      lastName: string
      emailAddress: string
      address: {
         streetAddr: string
         apartment: string?
         city: string
         state: string
         country: string?
      }?
      userMetadata: {
         createdTimestamp: number
         lastLoginIP: string?
         loginHistory: [any]?
         userProfileData: any?
      }
   }
```

Here is an unconventional JSTN text that mixes different object delineators and
uses unorthodox whitespace:

```
   {author:string;works:[{
     title:string
     year:     number?;
     classic:boolean;}]}
```

Here are several small JSTN texts:

```
   string

   number?

   boolean

   null

   [number]

   [string?]?
```

## References

### Normative References

```
   [RFC2119] Bradner, S., "Key words for use in RFCs to Indicate
             Requirement Levels", BCP 14, RFC 2119, March 1997.

   [RFC4234] Crocker, D. and P.  Overell, "Augmented BNF for Syntax
             Specifications: ABNF", RFC 4234, October 2005.

   [RFC7159]  Bray, T., Ed., "The JavaScript Object Notation (JSON) Data
              Interchange Format", RFC 7159, March 2014,.
```

## Author

```
	Tyler Christensen
	tyler9xp@gmail.com
```