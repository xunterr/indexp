package tokenizer

import (
	"slices"
	"testing"
)

func TestScanToken(t *testing.T) {
	tests := []struct {
		Name  string
		Input []byte
		Want  string
	}{
		{Name: "Parse simple string", Input: []byte("hello"), Want: "hello"},
		{Name: "Parse complicated string", Input: []byte("hello!@12a,"), Want: "hello"},
	}

	for _, test := range tests {
		tokenizer := NewTokenizer(test.Input)
		token, err := tokenizer.ScanToken()

		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
			continue
		}

		if token != test.Want {
			t.Errorf("Unexpected token. Have: %s, want: %s", token, test.Want)
			continue
		}
	}
}

func TestScanAll(t *testing.T) {
	tests := []struct {
		Name  string
		Input []byte
		Want  []string
	}{
		{Name: "Parse simple string", Input: []byte("lorem ipsum"), Want: []string{"lorem", "ipsum"}},
		{Name: "Parse complicated string", Input: []byte("hello!@12a,"), Want: []string{"hello", "!", "@", "12a", ","}},
	}

	for _, test := range tests {
		tokenizer := NewTokenizer(test.Input)
		tokens, err := tokenizer.ScanAll()

		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
			continue
		}

		if slices.Compare[[]string](tokens, test.Want) != 0 {
			t.Errorf("Unexpected token. Have: %v, want: %v", tokens, test.Want)
			continue
		}
	}
}
