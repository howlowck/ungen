package main

import (
	"testing"

	require "github.com/alecthomas/assert/v2"
)

func TestLexer(t *testing.T) {
	cases := []string{
		`// UNGEN: replace "World" with var.app_name`,
		`// UNGEN: delete 3 lines`,
		// TODO: still need to parse out braces and value
		`// UNGEN: replace "World" with "New {var.app_name}"`,
		`// UNGEN: replace "World" with kebabCase(var.app_name)`,
		`// UNGEN: replace "World" with titleCase(var.app_name)`,
		`// UNGEN:v1 if var.app_name then delete 3 lines`,
		`// UNGEN: if var.app_name then delete 3 lines else delete 1 line`,
		`# UNGEN: delete 1 line`,
		`// UNGEN:v1 replace "test" with var.appTest`,
		`// UNGEN: delete folder`,
		`// UNGEN: delete file`,
	}

	for _, c := range cases {
		_, err := Parse(c)
		// repr.Println(prog1)
		require.NoError(t, err)
	}
}
