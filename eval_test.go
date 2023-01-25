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

// UNGEN: delete 2 lines
REAL=false
REMOVE=false

// UNGEN: if var.use_dev then delete 1 line
DEBUG=true
	`

	vars := make(map[string]string)
	vars["app_port"] = "8000"
	vars["use_dev"] = "true"

	lines := regexp.MustCompile("\r?\n").Split(fileContent, -1)

	prog1, _ := Parse(`// UNGEN: replace "world" with "test"`)

	patch1 := prog1.Evaluate(lines, vars, 2)
	expected1 := []ContentPatch{{
		PatchType:     PatchReplace,
		OldLineNumber: 3,
		OldLineCount:  1,
		NewContent:    []string{"APP_NAME=hello-test"},
	}}

	eq1 := reflect.DeepEqual(patch1, expected1)
	if !eq1 {
		t.Error("Failed patch1 does not equal expected", eq1)
	}

	prog2, _ := Parse(`// UNGEN: replace "3000" with var.app_port`)
	patch2 := prog2.Evaluate(lines, vars, 5)
	expected2 := []ContentPatch{{
		PatchType:     PatchReplace,
		OldLineNumber: 6,
		OldLineCount:  1,
		NewContent:    []string{"APP_PORT=8000"},
	}}

	eq2 := reflect.DeepEqual(patch2, expected2)
	if !eq2 {
		t.Error("Failed patch2 does not equal expected", eq2)
	}

	prog3, _ := Parse(`// UNGEN: delete 2 lines`)
	patch3 := prog3.Evaluate(lines, vars, 8)
	expected3 := []ContentPatch{{
		PatchType:     PatchDelete,
		OldLineNumber: 9,
		OldLineCount:  2,
		NewContent:    []string{},
	}}

	eq3 := reflect.DeepEqual(patch3, expected3)
	if !eq3 {
		t.Error("Failed patch3 does not equal expected", eq3)
	}
}
