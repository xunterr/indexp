package indexer

import (
	"io"
	"strings"
)

type Predicate func(b byte) bool

type Tokenizer struct {
	source  []byte
	start   int
	current int
}

func NewTokenizer(source []byte) *Tokenizer {
	return &Tokenizer{
		source:  source,
		start:   0,
		current: 0,
	}
}

func (t *Tokenizer) ScanToken() (string, error) {
	if t.isAtEnd() {
		return "", io.EOF
	}
	t.start = t.current
	c := t.source[t.current]
	t.current++
	result := ""
	switch c {
	case '\n':
	case ' ':
	case '\r':
	case '\t':
		break
	default:
		if isDigit(c) {
			result = t.scanWhile(isDigit)
		} else if isAlfa(c) {
			result = t.scanWhile(isAlfa)
		}
	}
	return strings.ToLower(result), nil
}

func (t *Tokenizer) ScanAll() (tokens []string, err error) {
	for {
		token, scanErr := t.ScanToken()
		if scanErr != nil {
			if scanErr == io.EOF {
				break
			}
			err = scanErr
			return
		}
		if token != "" {
			tokens = append(tokens, token)
		}
	}
	return
}

func (t *Tokenizer) scanWhile(predicate Predicate) string {
	for predicate(t.peek()) && !t.isAtEnd() {
		t.current++
	}
	res := string(t.source[t.start:t.current])

	return res
}

func (t *Tokenizer) peek() byte {
	if t.isAtEnd() {
		return '\000'
	}
	return t.source[t.current]
}

func (t *Tokenizer) peekNext() byte {
	if t.current+1 > len(t.source) {
		return '\000'
	}
	return t.source[t.current+1]
}

func (t Tokenizer) isAtEnd() bool {
	return t.current >= len(t.source)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlfa(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z')
}
