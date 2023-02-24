package main

import (
	"fmt"
	"reflect"
	"testing"
)

type PatchTestCase struct {
	title    string
	original []string
	patch    ContentPatch
	expected []string
}

func TestContentPatch(t *testing.T) {
	testCases := []PatchTestCase{
		{
			"PatchReplace",
			[]string{
				`// this is a ini file`,
				`// UNGEN: replace "world" with "test"`,
				`APP_NAME=hello-world`,
				``,
				`// UNGEN: replace "3000" with var.app_port`,
				`APP_PORT=3000`,
				``,
				`// UNGEN: delete 2 lines`,
				`REAL=false`,
				`REMOVE=false`,
				``,
				`// UNGEN: if var.use_dev then delete 1 line`,
				`DEBUG=true`,
			},
			ContentPatch{
				PatchType:     PatchReplace,
				OldLineNumber: 2,
				OldLineCount:  2,
				NewContent:    []string{"APP_NAME=hello-test"},
			},
			[]string{
				`// this is a ini file`,
				`APP_NAME=hello-test`,
				``,
				`// UNGEN: replace "3000" with var.app_port`,
				`APP_PORT=3000`,
				``,
				`// UNGEN: delete 2 lines`,
				`REAL=false`,
				`REMOVE=false`,
				``,
				`// UNGEN: if var.use_dev then delete 1 line`,
				`DEBUG=true`,
			},
		},
		{
			"PatchDelete",
			[]string{
				`// this is a ini file`,
				`// UNGEN: replace "world" with "test"`,
				`APP_NAME=hello-world`,
				``,
				`// UNGEN: replace "3000" with var.app_port`,
				`APP_PORT=3000`,
				``,
				`// UNGEN: delete 2 lines`,
				`REAL=false`,
				`REMOVE=false`,
				``,
				`// UNGEN: if var.use_dev then delete 1 line`,
				`DEBUG=true`,
			},
			ContentPatch{
				PatchType:     PatchDelete,
				OldLineNumber: 2,
				OldLineCount:  2,
				NewContent:    []string{},
			},
			[]string{
				`// this is a ini file`,
				``,
				`// UNGEN: replace "3000" with var.app_port`,
				`APP_PORT=3000`,
				``,
				`// UNGEN: delete 2 lines`,
				`REAL=false`,
				`REMOVE=false`,
				``,
				`// UNGEN: if var.use_dev then delete 1 line`,
				`DEBUG=true`,
			},
		},
		{
			"PatchReplace (in the form of Insert)",
			[]string{
				`// this is a ini file`,
				`// UNGEN: replace "world" with "test"`,
				`APP_NAME=hello-world`,
				``,
				`// UNGEN: replace "3000" with var.app_port`,
				`APP_PORT=3000`,
				`// UNGEN: insert var.line`,
				``,
				`// UNGEN: delete 2 lines`,
				`REAL=false`,
				`REMOVE=false`,
				``,
				`// UNGEN: if var.use_dev then delete 1 line`,
				`DEBUG=true`,
			},
			ContentPatch{
				PatchType:     PatchReplace,
				OldLineNumber: 8,
				OldLineCount:  0,
				NewContent: []string{
					"LINE1=1",
					"LINE2=2",
				},
			},
			[]string{
				`// this is a ini file`,
				`// UNGEN: replace "world" with "test"`,
				`APP_NAME=hello-world`,
				``,
				`// UNGEN: replace "3000" with var.app_port`,
				`APP_PORT=3000`,
				`// UNGEN: insert var.line`,
				"LINE1=1",
				"LINE2=2",
				``,
				`// UNGEN: delete 2 lines`,
				`REAL=false`,
				`REMOVE=false`,
				``,
				`// UNGEN: if var.use_dev then delete 1 line`,
				`DEBUG=true`,
			},
		},
		{
			"PatchReplace (in the form of removing command in insert)",
			[]string{
				`// this is a ini file`,
				`// UNGEN: replace "world" with "test"`,
				`APP_NAME=hello-world`,
				``,
				`// UNGEN: replace "3000" with var.app_port`,
				`APP_PORT=3000`,
				`// UNGEN: insert var.line`,
				``,
				`// UNGEN: delete 2 lines`,
				`REAL=false`,
				`REMOVE=false`,
				``,
				`// UNGEN: if var.use_dev then delete 1 line`,
				`DEBUG=true`,
			},
			ContentPatch{
				PatchType:     PatchReplace,
				OldLineNumber: 7,
				OldLineCount:  1,
				NewContent: []string{
					"LINE1=1",
					"LINE2=2",
				},
			},
			[]string{
				`// this is a ini file`,
				`// UNGEN: replace "world" with "test"`,
				`APP_NAME=hello-world`,
				``,
				`// UNGEN: replace "3000" with var.app_port`,
				`APP_PORT=3000`,
				"LINE1=1",
				"LINE2=2",
				``,
				`// UNGEN: delete 2 lines`,
				`REAL=false`,
				`REMOVE=false`,
				``,
				`// UNGEN: if var.use_dev then delete 1 line`,
				`DEBUG=true`,
			},
		},
	}

	for i, c := range testCases {
		actual := c.patch.Apply(c.original)
		eq := reflect.DeepEqual(actual, c.expected)
		if !eq {
			t.Error("Failed patch does not equal expected. Test Case index:", i)
		}
		fmt.Println("✅ Test Case:", c.title)
	}
}
