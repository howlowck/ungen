package main

import (
	"reflect"
	"regexp"
	"testing"
)

type EvalTestCase struct {
	Context  EvalContext
	Command  string
	Expected []Patch
}

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

// UNGEN: replace "false" with titleCase(var.use_dev)
TITLE=false
	`

	vars := make(map[string]string)
	vars["app_port"] = "8000"
	vars["use_dev"] = "true"

	lines := regexp.MustCompile("\r?\n").Split(fileContent, -1)

	testCases := []EvalTestCase{
		{
			Context: EvalContext{
				lines:             lines,
				path:              ".env.test",
				vars:              vars,
				programLineNumber: 2,
			},
			Command: `// UNGEN: replace "world" with "test"`,
			Expected: []Patch{{
				Content: &ContentPatch{
					PatchType:     PatchReplace,
					OldLineNumber: 3,
					OldLineCount:  1,
					NewContent:    []string{"APP_NAME=hello-test"},
				},
			}},
		},
		{
			Context: EvalContext{
				lines:             lines,
				path:              ".env.test",
				vars:              vars,
				programLineNumber: 5,
			},
			Command: `// UNGEN: replace "3000" with var.app_port`,
			Expected: []Patch{{
				Content: &ContentPatch{
					PatchType:     PatchReplace,
					OldLineNumber: 6,
					OldLineCount:  1,
					NewContent:    []string{"APP_PORT=8000"},
				},
			}},
		},
		{
			Context: EvalContext{
				lines:             lines,
				path:              ".env.test",
				vars:              vars,
				programLineNumber: 8,
			},
			Command: `// UNGEN: delete 2 lines`,
			Expected: []Patch{{
				Content: &ContentPatch{
					PatchType:     PatchDelete,
					OldLineNumber: 9,
					OldLineCount:  2,
					NewContent:    []string{},
				},
			}},
		},
		// {
		// 	Context: EvalContext{
		// 		lines:             lines,
		// 		path:              ".env.test",
		// 		vars:              vars,
		// 		programLineNumber: 15,
		// 	},
		// 	Command: `// UNGEN: replace "false" with titleCase(var.use_dev)`,
		// 	Expected: []Patch{{
		// 		Content: &ContentPatch{
		// 			PatchType:     PatchReplace,
		// 			OldLineNumber: 15,
		// 			OldLineCount:  2,
		// 			NewContent:    []string{"TITLE=True"},
		// 		},
		// 	}},
		// },
	}

	for i, c := range testCases {
		p, _ := Parse(c.Command)
		actual := p.Evaluate(c.Context)
		eq := reflect.DeepEqual(actual, c.Expected)
		if !eq {
			t.Error("Failed actual does not equal expected at index:", i)
		}
	}
}
