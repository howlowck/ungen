package internal

import (
	"testing"

	require "github.com/alecthomas/assert/v2"
	"github.com/alecthomas/repr"
)

func TestLexerLite(t *testing.T) {
	cases := []string{
		`UNGEN:gpt @ I want you to replace ${ five } with ${ kebabCase(var.app_name) } @`,
	}

	for _, c := range cases {
		p, err := LiteParse(c)
		repr.Println(c, p)
		require.NoError(t, err)
	}
}
