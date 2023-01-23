package main

import (
	"testing"

	require "github.com/alecthomas/assert/v2"
	"github.com/alecthomas/repr"
)

func TestLexer(t *testing.T) {
	prog1, err := Parse(`// UNGEN: replace "World" with var.app_name`)
	require.NoError(t, err)
	repr.Println(prog1)
	prog2, err := Parse(`// UNGEN: delete 3 lines`)
	require.NoError(t, err)
	repr.Println(prog2)
}
