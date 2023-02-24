package main

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

type GatherTestCase struct {
	Context       Context
	ExpectedKey   string
	ExpectedValue []string
}

func lines(s string) []string {
	return regexp.MustCompile("\r?\n").Split(s, -1)
}

func TestGather(t *testing.T) {
	vars := make(map[string]string)
	vars["useTypescript"] = "true"
	testCases := []GatherTestCase{
		{
			Context: Context{
				lines: lines(`FILE=.env.test
APP_NAME=hello-test
// UNGEN: cut ln 1-2 to cb.message
				`),
				path:              ".env.test",
				vars:              vars,
				clipboard:         make(map[string][]string),
				keepLine:          true,
				programLineNumber: 3,
			},
			ExpectedKey: "message",
			ExpectedValue: []string{
				"FILE=.env.test",
				"APP_NAME=hello-test",
			},
		},
		{
			Context: Context{
				lines: lines(`FILE=.env.test
APP_NAME=hello-test
// UNGEN: copy ln 1 to cb.message
						`),
				path:              ".env.test",
				vars:              vars,
				clipboard:         make(map[string][]string),
				keepLine:          true,
				programLineNumber: 3,
			},
			ExpectedKey: "message",
			ExpectedValue: []string{
				"FILE=.env.test",
			},
		},
		{
			Context: Context{
				lines: lines(`FILE=.env.test
APP_NAME=hello-test
// UNGEN: copy next 2 lines to cb.message
SOME_VAR=hello
SOME_OTHER_VAR=world
						`),
				path:              ".env.test",
				vars:              vars,
				clipboard:         make(map[string][]string),
				keepLine:          true,
				programLineNumber: 3,
			},
			ExpectedKey: "message",
			ExpectedValue: []string{
				"SOME_VAR=hello",
				"SOME_OTHER_VAR=world",
			},
		},
		{
			Context: Context{
				lines: lines(`FILE=.env.test
APP_NAME=hello-test
// UNGEN: cut next 2 lines to cb.message
SOME_VAR=hello
SOME_OTHER_VAR=world
						`),
				path:              ".env.test",
				vars:              vars,
				clipboard:         make(map[string][]string),
				keepLine:          true,
				programLineNumber: 3,
			},
			ExpectedKey: "message",
			ExpectedValue: []string{
				"SOME_VAR=hello",
				"SOME_OTHER_VAR=world",
			},
		},
		{
			Context: Context{
				lines: lines(`FILE=.env.test
APP_NAME=hello-test
// UNGEN: if var.useTypescript then cut next 2 lines to cb.message
SOME_VAR=hello
SOME_OTHER_VAR=world
						`),
				path:              ".env.test",
				vars:              vars,
				clipboard:         make(map[string][]string),
				keepLine:          true,
				programLineNumber: 3,
			},
			ExpectedKey: "message",
			ExpectedValue: []string{
				"SOME_VAR=hello",
				"SOME_OTHER_VAR=world",
			},
		},
		{
			Context: Context{
				lines: lines(`FILE=.env.test
APP_NAME=hello-test
// UNGEN: if "false" then cut next 1 line to cb.message
SOME_VAR=hello
SOME_OTHER_VAR=world
						`),
				path:              ".env.test",
				vars:              vars,
				clipboard:         make(map[string][]string),
				keepLine:          true,
				programLineNumber: 3,
			},
			ExpectedKey:   "message",
			ExpectedValue: nil,
		},
	}

	for i, c := range testCases {
		cmd := c.Context.lines[c.Context.programLineNumber-1]
		p, _ := Parse(cmd)
		p.Gather(&c.Context)
		actual := c.Context.clipboard[c.ExpectedKey]

		eq := reflect.DeepEqual(actual, c.ExpectedValue)

		if !eq {
			t.Errorf("Test case %d failed. Expected %v, got %v", i, c.ExpectedValue, actual)
		}
		fmt.Println("✅", cmd)
	}
}
