package main

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

type EvalTestCase struct {
	Context  Context
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

// UNGEN: copy next 1 line to cb.description
DESCRIPTION=changeme

// UNGEN: cut ln 23-24 to cb.description2
DESCRIPTION1=changeme1
DESCRIPTION2=changeme2

// UNGEN: copy ln 1 to cb.description3

// UNGEN: copy ln 5-7 to cb.description4

// UNGEN: insert cb.description
`

	vars := make(map[string]string)
	vars["app_port"] = "8000"
	vars["use_dev"] = "true"
	vars["app_name"] = "test-app"

	clipboard := make(map[string][]string)

	lines := regexp.MustCompile("\r?\n").Split(fileContent, -1)

	testCases := []EvalTestCase{
		{
			Context: Context{
				lines:             lines,
				path:              ".env.test",
				vars:              vars,
				clipboard:         clipboard,
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
			Context: Context{
				lines:             lines,
				path:              ".env.test",
				vars:              vars,
				clipboard:         clipboard,
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
			Context: Context{
				lines:             lines,
				path:              ".env.test",
				vars:              vars,
				clipboard:         clipboard,
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
			Context: Context{
				lines:             lines,
				path:              ".env.test",
				vars:              vars,
				clipboard:         clipboard,
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
			Context: Context{
				lines:             lines,
				path:              ".env.test",
				vars:              vars,
				clipboard:         clipboard,
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
			Context: Context{
				lines:             lines,
				path:              ".env.test",
				vars:              vars,
				clipboard:         clipboard,
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
			Context: Context{
				lines:             lines,
				path:              ".env.test",
				vars:              vars,
				clipboard:         clipboard,
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
			Context: Context{
				lines:             lines,
				path:              "test/.env.test",
				vars:              vars,
				clipboard:         clipboard,
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
		{
			Context: Context{
				lines:             lines,
				path:              "test/.env.test",
				vars:              vars,
				clipboard:         clipboard,
				keepLine:          true,
				programLineNumber: 27,
			},
			Command: lines[26], // copy
			Expected: []Patch{{
				Content: &ContentPatch{
					PatchType:     PatchDelete,
					OldLineNumber: 28,
					OldLineCount:  0,
					NewContent:    []string{},
				},
			}},
		},
		{
			Context: Context{
				lines:             lines,
				path:              "test/.env.test",
				vars:              vars,
				clipboard:         clipboard,
				keepLine:          false,
				programLineNumber: 30,
			},
			Command: lines[29], // cut
			Expected: []Patch{{
				Content: &ContentPatch{
					PatchType:     PatchDelete,
					OldLineNumber: 23,
					OldLineCount:  2,
					NewContent:    []string{},
				},
			}, {
				Content: &ContentPatch{
					PatchType:     PatchDelete,
					OldLineNumber: 30,
					OldLineCount:  1,
					NewContent:    []string{},
				},
			}},
		},
		{
			Context: Context{
				lines:             lines,
				path:              "test/.env.test",
				vars:              vars,
				clipboard:         clipboard,
				keepLine:          false,
				programLineNumber: 34,
			},
			Command: lines[33], // copy ln 1
			Expected: []Patch{{
				Content: &ContentPatch{
					PatchType:     PatchDelete,
					OldLineNumber: 34,
					OldLineCount:  1,
					NewContent:    []string{},
				},
			}},
		},
		{
			Context: Context{
				lines:             lines,
				path:              "test/.env.test",
				vars:              vars,
				clipboard:         clipboard,
				keepLine:          false,
				programLineNumber: 36,
			},
			Command: lines[35], // copy ln 5-7
			Expected: []Patch{{
				Content: &ContentPatch{
					PatchType:     PatchDelete,
					OldLineNumber: 36,
					OldLineCount:  1,
					NewContent:    []string{},
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
