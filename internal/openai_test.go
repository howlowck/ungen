package internal

import (
	"os"
	"regexp"
	"testing"
)

type OpenAITestCase struct {
	Context  Context
	Command  string
	Expected []Patch
}

func TestOpenAI(t *testing.T) {
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

// UNGEN: if var.app_name == "test-app" then replace "false" with "true" else delete file
TEST=false
`
	apiKey := os.Getenv("UNGEN_OPENAI_API_KEY")

	vars := make(map[string]string)
	vars["app_port"] = "8000"
	vars["use_dev"] = "true"
	vars["app_name"] = "test-app"

	clipboard := make(map[string][]string)

	lines := regexp.MustCompile("\r?\n").Split(fileContent, -1)

	testCases := []OpenAITestCase{
		{
			Context: Context{
				Lines:             lines,
				Path:              ".env.test",
				Vars:              vars,
				Clipboard:         clipboard,
				KeepLine:          true,
				ProgramLineNumber: 2,
			},
			Command: "change APP_NAME to kebabCase(var.app_name) on line 10",
			Expected: []Patch{{
				Content: &ContentPatch{
					PatchType:     PatchReplace,
					OldLineNumber: 3,
					OldLineCount:  1,
					NewContent:    []string{"APP_NAME=hello-test"},
				},
			}},
		},
	}

	for _, c := range testCases {

		OpenAIEval(apiKey, c.Context, c.Command)
		// for _, ap := range actual {
		// 	if ap.File != nil {
		// 		fmt.Println(*ap.File)
		// 	} else {
		// 		fmt.Println(*ap.Content)
		// 	}
		// }
		// eq := reflect.DeepEqual(actual, c.Expected)
		// if !eq {
		// 	t.Error("Failed actual does not equal expected at index:", i)
		// }
	}
}
