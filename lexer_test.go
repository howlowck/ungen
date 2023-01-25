package main

import (
	"testing"

	require "github.com/alecthomas/assert/v2"
)

func TestLexer(t *testing.T) {
	cases := []string{
		`// UNGEN: replace "World" with var.app_name`,
		`// UNGEN: delete 3 lines`,
		`// UNGEN: if var.app_name then delete 3 lines`,
		`// UNGEN: delete file`,
	}

	for _, c := range cases {
		_, err := Parse(c)
		// repr.Println(prog1)
		require.NoError(t, err)
	}
}
