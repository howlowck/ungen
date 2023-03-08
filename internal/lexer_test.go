package internal

import (
	"testing"

	require "github.com/alecthomas/assert/v2"
	"github.com/alecthomas/repr"
)

func TestLexer(t *testing.T) {
	cases := []string{
		`UNGEN: replace "World" with var.app_name`,
		`UNGEN: delete 3 lines`,
		`UNGEN: replace "World" with "New ${var.app_name}"`, // NOT IMPLEMENTED YET
		`UNGEN: replace "World" with kebabCase(var.app_name)`,
		`UNGEN: replace "World" with substitute(var.app_name, "-", "")`,
		`UNGEN:v1 if var.app_name then delete 3 lines`,
		`UNGEN: if var.app_name then delete 3 lines else delete 1 line`,
		`UNGEN: delete 1 line`,
		`UNGEN:v1 replace "test" with var.appTest`,
		`UNGEN: delete folder`,
		`UNGEN: delete file`,
		`UNGEN: copy next 1 line to cb.keyVault`,
		`UNGEN: cut next 1 line to cb.keyVault`,
		`UNGEN: copy ln 10 to cb.keyVault`,
		`UNGEN: copy ln 10-12 to cb.keyVault`,
		`UNGEN: insert cb.keyVault`,
		`UNGEN: if var.useTypescript == "yes" then insert cb.tsconfig else delete file`,
		`UNGEN: inject file:package.json on ln 5 'replace "simple-app" with var.appName'`,
		`UNGEN:scoped inject file:package.json on ln 5 'replace "simple-app" with var.appName'`,
	}

	for _, c := range cases {
		p, err := Parse(c)
		repr.Println(c, p)
		require.NoError(t, err)
	}
}
