package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Parser struct {
	input string
	pos   int
}

func NewParser(input string) *Parser {
	return &Parser{input: input, pos: 0}
}

func (P *Parser) Parse() (any, error) {
	P.skipWhitespace()
	return P.parseValue()
}

func (P *Parser) skipWhitespace() {
	for P.pos < len(P.input) && unicode.IsSpace(rune(P.input[P.pos])) {
		P.pos++
	}
}

func (P *Parser) parseObject() (any, error) {
	P.skipWhitespace()
	res := make(map[string]any)
	if P.input[P.pos] != '{' {
		return nil, fmt.Errorf("expcted { found at %d", P.pos)
	}
	P.pos++

	P.skipWhitespace()
	for len(P.input) > P.pos && P.input[P.pos] != '}' {
		P.skipWhitespace()
		key, err := P.parseString()

		if err != nil {
			return nil, err
		}
		P.skipWhitespace()
		if len(P.input) <= P.pos {
			break
		}
		if P.input[P.pos] == ':' {
			P.pos++
		}
		P.skipWhitespace()
		val, err := P.parseValue()
		res[key] = val
		P.skipWhitespace()
		if P.input[P.pos] == ',' {
			P.pos++
		} else if P.input[P.pos] == '}' {
			P.pos++
			break
		} else {
			return nil, fmt.Errorf("expected , or } at pos %d", P.pos)
		}
	}

	return res, nil
}

func (P *Parser) parseValue() (any, error) {
	if P.pos >= len(P.input) {
		return nil, fmt.Errorf("unexpected EOF")
	}

	switch P.input[P.pos] {
	case '{':
		return P.parseObject()
	case '[':
		return P.parseArray()
	case '"':
		return P.parseString()
	case 't', 'f':
		return P.parseBoolean()
	case 'n':
		return P.parseNull()
	default:
		if unicode.IsDigit(rune(P.input[P.pos])) || P.input[P.pos] == '-' {
			return P.parseNumber()
		}
		return nil, fmt.Errorf("unexpected char: %c at pos %d", P.input[P.pos], P.pos)
	}
}

func (P *Parser) parseNull() (any, error) {
	if P.pos >= len(P.input) {
		return nil, nil
	}
	P.skipWhitespace()
	if P.pos+3 >= len(P.input) || P.input[P.pos:P.pos+4] != "null" {
		return nil, fmt.Errorf("expected null found some weird stuff only god knows")
	}

	P.pos += 4
	return nil, nil
}

func (P *Parser) parseArray() ([]any, error) {
	if P.pos >= len(P.input) {
		return nil, nil
	}
	var arr []any

	P.pos++
	for P.pos < len(P.input) && P.input[P.pos] != ']' {
		P.skipWhitespace()
		val, err := P.Parse()
		if err != nil {
			return nil, err
		}

		arr = append(arr, val)
		if P.input[P.pos] == ',' {
			P.pos++
		}
	}
	if P.input[P.pos] == ']' {
		fmt.Println("yes")
		P.pos++
	}
	return arr, nil
}

func (P *Parser) parseNumber() (any, error) {
	if P.pos >= len(P.input) {
		return nil, nil
	}
	if P.pos >= len(P.input) {
		return nil, nil
	}
	P.skipWhitespace()
	numStr := ""
	for P.pos < len(P.input) && unicode.IsDigit(rune(P.input[P.pos])) {
		numStr += string(P.input[P.pos])
		P.pos++
	}

	if numStr == "" {
		return nil, fmt.Errorf("epected number at pos %d", P.pos)
	}
	return strconv.ParseInt(numStr, 10, 32)
}

func (P *Parser) parseBoolean() (bool, error) {
	if P.pos >= len(P.input) {
		return false, nil
	}
	P.skipWhitespace()
	if P.input[P.pos] == 't' {
		P.pos += 4
		return true, nil
	}
	P.pos += 5
	return false, nil
}

func (P *Parser) parseString() (string, error) {
	if P.pos >= len(P.input) {
		return "", fmt.Errorf("unexpected eof")
	}
	P.skipWhitespace()
	str := ""
	P.pos++

	for P.input[P.pos] != '"' {
		str += string(P.input[P.pos])
		P.pos++
	}
	P.pos++
	return str, nil
}

func main() {
	inputs := []string{
		`{"aha":["meow","kehe"], "kaha":"yahan"}`,
		`{"emptyArray":[], "emptyObj":{}}`,
		`{"num":123, "bool":true, "nullVal":null}`,
		`[1, 2, 3, 4, 5]`,
		`["string", 42,  null, {"nested":"obj"}, [1,2]]`,
		`[{"nested":"obj"}]`,
		`"just a string"`,
		`12345`,
		`true`,
		`false`,
		`null`,
		`{"nestedArr":[{"a":1},{"b":2}], "nestedObj":{"x":10,"y":20}}`,
		`[]`,
		`{}`,
		`[[[[]]]]`,
		`{"a":{"b":{"c":null}}}`,
	}
	for _, input := range inputs {
		parser := NewParser(strings.TrimSpace(input))
		res, err := parser.Parse()
		if err != nil {
			fmt.Println("error parsing: ", err)
			return
		}
		fmt.Println(res)
	}
}
