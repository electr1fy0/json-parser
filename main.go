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
	res := make(map[string]any)
	P.pos++
	var key string
	var val any
	var err error
	for P.input[P.pos] != ':' {
		key, err = P.parseString()
		if err != nil {
			return nil, err
		}
	}
	P.pos++
	for P.input[P.pos] != '}' {
		val, err = P.parseValue()
		if err != nil {
			return nil, err
		}

	}

	res[key] = val
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
	if P.pos+3 >= len(P.input) || P.input[P.pos:P.pos+4] != "null" {
		return nil, fmt.Errorf("expected null found some weird stuff only god knows")
	}

	P.pos += 4
	return nil, nil
}

func (P *Parser) parseArray() ([]any, error) {
	var arr []any

	P.pos++
	for P.input[P.pos] != ']' {
		val, err := P.Parse()
		if err != nil {
			return nil, err
		}

		arr = append(arr, val)
		if P.input[P.pos] == ',' {
			P.pos++
		}
	}

	return arr, nil
}

func (P *Parser) parseNumber() (any, error) {
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
	if P.input[P.pos] == 't' {
		P.pos++
		return true, nil
	}
	P.pos++
	return false, nil
}

func (P *Parser) parseString() (string, error) {
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
	inputs := []string{`[1,2,3,4,5]`, `{"meow":"meow2"}`, `323244`, "null"}

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
