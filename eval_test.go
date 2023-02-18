package main

import (
	"fmt"
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

// UNGEN: replace "false" with upperCase(var.use_dev)
TITLE=false

// UNGEN: replace "changeme" with substitute(var.app_name, "-", "")
STORAGE_ACCOUNT=changeme

// UNGEN: replace "changeme" with concat(var.app_name, " ", "welcomes you")
STARTUP_MESSAGE=changeme

// UNGEN: delete file
// UNGEN: delete folder
`

	vars := make(map[string]string)
	vars["app_port"] = "8000"
	vars["use_dev"] = "true"
	vars["app_name"] = "test-app"

	lines := regexp.MustCompile("\r?\n").Split(fileContent, -1)

	testCases := []EvalTestCase{
		{
			Context: EvalContext{
				lines:             lines,
				path:              ".env.test",
				vars:              vars,
				keepLine:          true,
				programLineNumber: 2,
			},
			Command: lines[1],
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
				keepLine:          true,
				programLineNumber: 5,
			},
			Command: lines[4],
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
				keepLine:          true,
				programLineNumber: 8,
			},
			Command: lines[7],
			Expected: []Patch{{
				Content: &ContentPatch{
					PatchType:     PatchDelete,
					OldLineNumber: 9,
					OldLineCount:  2,
					NewContent:    []string{},
				},
			}},
		},
		{
			Context: EvalContext{
				lines:             lines,
				path:              ".env.test",
				vars:              vars,
				keepLine:          true,
				programLineNumber: 15,
			},
			Command: lines[14],
			Expected: []Patch{{
				Content: &ContentPatch{
					PatchType:     PatchReplace,
					OldLineNumber: 16,
					OldLineCount:  1,
					NewContent:    []string{"TITLE=TRUE"},
				},
			}},
		},
		{
			Context: EvalContext{
				lines:             lines,
				path:              ".env.test",
				vars:              vars,
				keepLine:          true,
				programLineNumber: 18,
			},
			Command: lines[17],
			Expected: []Patch{{
				Content: &ContentPatch{
					PatchType:     PatchReplace,
					OldLineNumber: 19,
					OldLineCount:  1,
					NewContent:    []string{"STORAGE_ACCOUNT=testapp"},
				},
			}},
		},
		{
			Context: EvalContext{
				lines:             lines,
				path:              ".env.test",
				vars:              vars,
				keepLine:          true,
				programLineNumber: 21,
			},
			Command: lines[20],
			Expected: []Patch{{
				Content: &ContentPatch{
					PatchType:     PatchReplace,
					OldLineNumber: 22,
					OldLineCount:  1,
					NewContent:    []string{"STARTUP_MESSAGE=test-app welcomes you"},
				},
			}},
		},
		{
			Context: EvalContext{
				lines:             lines,
				path:              ".env.test",
				vars:              vars,
				keepLine:          true,
				programLineNumber: 24,
			},
			Command: lines[23],
			Expected: []Patch{{
				File: &FilePatch{
					FileOp:     FileDelete,
					TargetPath: ".env.test",
				},
			}},
		},
		{
			Context: EvalContext{
				lines:             lines,
				path:              "test/.env.test",
				vars:              vars,
				keepLine:          true,
				programLineNumber: 25,
			},
			Command: lines[24],
			Expected: []Patch{{
				File: &FilePatch{
					FileOp:     DirectoryDelete,
					TargetPath: "test/",
				},
			}},
		},
	}

	for i, c := range testCases {
		p, _ := Parse(c.Command)
		actual := p.Evaluate(c.Context)
		for _, ap := range actual {
			if ap.File != nil {
				fmt.Println(*ap.File)
			} else {
				fmt.Println(*ap.Content)
			}
		}
		eq := reflect.DeepEqual(actual, c.Expected)
		if !eq {
			t.Error("Failed actual does not equal expected at index:", i)
		}
	}
}
