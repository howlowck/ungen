package main

import (
	"reflect"
	"testing"
)

type PatchTestCase struct {
	original []string
	patch    ContentPatch
	expected []string
}

func TestPatch(t *testing.T) {
	testCases := []PatchTestCase{
		{
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
	}

	for i, c := range testCases {
		actual := c.patch.Apply(c.original)
		eq := reflect.DeepEqual(actual, c.expected)
		// fmt.Println(actual, c.expected)
		if !eq {
			t.Error("Failed patch does not equal expected. Test Case index:", i)
		}
	}
}
