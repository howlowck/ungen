package main

import (
	"reflect"
	"regexp"
	"testing"
)

func TestEval(t *testing.T) {

	fileContent := `// this is a ini file
// UNGEN: replace "world" with "test"
APP_NAME=hello-world

// UNGEN: replace "3000" with var.app_port
APP_PORT=3000

// UNGEN: if var.use_dev then delete 1 line
DEBUG=true

// UNGEN: delete 2 lines
REAL=false
REMOVE=false
	`

	vars := make(map[string]string)
	vars["app_port"] = "8000"
	vars["use_dev"] = "true"

	prog1, _ := Parse(`// UNGEN: replace "world" with "test"`)
	lines := regexp.MustCompile("\r?\n").Split(fileContent, -1)

	patch1 := prog1.Evaluate(lines, vars, 2)
	expected1 := []Patch{{
		PatchType:     PatchReplace,
		OldLineNumber: 3,
		OldLineCount:  1,
		NewContent:    "APP_NAME=hello-test",
	}}

	eq1 := reflect.DeepEqual(patch1, expected1)
	if !eq1 {
		t.Error("Failed patch1 does not equal expected", eq1)
	}

	prog2, _ := Parse(`// UNGEN: replace "3000" with var.app_port`)
	patch2 := prog2.Evaluate(lines, vars, 5)
	expected2 := []Patch{{
		PatchType:     PatchReplace,
		OldLineNumber: 6,
		OldLineCount:  1,
		NewContent:    "APP_PORT=8000",
	}}

	eq2 := reflect.DeepEqual(patch2, expected2)
	if !eq2 {
		t.Error("Failed patch2 does not equal expected", eq2)
	}
	// prog3, _ := Parse(`// UNGEN: delete 2 lines`)

}
