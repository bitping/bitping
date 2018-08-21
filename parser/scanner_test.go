package parser_test

import (
	"strings"
	"testing"

	parser "github.com/auser/bitping/parser"
)

func TestScanner_Scan(t *testing.T) {
	var tests = []struct {
		s   string
		tok parser.Token
		lit string
	}{
		{s: ``, tok: parser.EOF},
		{s: `#`, tok: parser.ILLEGAL, lit: `#`},
		{s: ` `, tok: parser.WS, lit: ` `},
		{s: "\t", tok: parser.WS, lit: "\t"},
		{s: "\n", tok: parser.WS, lit: "\n"},

		// Identifiers
		{s: `0xBD6d79F3f02584cfcB754437Ac6776c4C6E0a0eC`, tok: parser.ADDRESS, lit: `0xBD6d79F3f02584cfcB754437Ac6776c4C6E0a0eC`},
		{s: `_x222`, tok: parser.ILLEGAL, lit: `_`},
		{s: `4`, tok: parser.DIGITS, lit: "4"},
		{s: `4e10`, tok: parser.DIGITS, lit: "4e10"},

		// Keywords
		{s: `SELECT`, tok: parser.SELECT, lit: "SELECT"},
		{s: `AND`, tok: parser.AND, lit: "AND"},
		{s: `<`, tok: parser.LESSTHAN, lit: "<"},
		{s: `>`, tok: parser.GREATERTHAN, lit: ">"},
		{s: "=", tok: parser.EQUAL, lit: "="},
		{s: "!", tok: parser.NOT, lit: "!"},
	}

	for i, tt := range tests {
		s := parser.NewScanner(strings.NewReader(tt.s))
		tok, lit := s.Scan()
		if tt.tok != tok {
			t.Errorf("%d, %q token mismatch: exp=%q, got=%q <%q>", i, tt.s, tt.tok, tok, lit)
		} else if tt.lit != lit {
			t.Errorf("%d. %q literal mismatch: exp=%q, got=%q", i, tt.s, tt.lit, lit)
		}
	}

}
