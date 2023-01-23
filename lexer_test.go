package main

import (
	"testing"

	require "github.com/alecthomas/assert/v2"
	"github.com/alecthomas/repr"
)

func TestLexer(t *testing.T) {
	prog1, err := Parse(`// UNGEN: replace "World" with var.app_name`)
	repr.Println(prog1)
	require.NoError(t, err)
	prog2, err := Parse(`// UNGEN: delete 3 lines`)
	repr.Println(prog2)
	require.NoError(t, err)
	prog3, err := Parse(`// UNGEN: if var.app_name then delete 3 lines`)
	repr.Println(prog3)
	require.NoError(t, err)
}
