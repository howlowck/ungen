package internal

import (
	"fmt"
	"testing"
)

type InjectTestCase struct {
	Context         InjectionContext
	UngenCmd        string
	ExpectedHistory map[string][]int
	ExpectedContent map[string][]string
}

func TestInject(t *testing.T) {
	testCases := []InjectTestCase{
		{
			Context: InjectionContext{
				DotFilePath:      "../examples/simple-nodejs/.ungen",
				InjectionHistory: map[string][]int{},
				InjectionContent: map[string][]string{
					"../examples/simple-nodejs/package.json": {
						`{`,
						`  "name": "simple-nodejs",`,
						`  "version": "1.0.0",`,
						`  "description": "",`,
						`  "main": "index.js",`,
						`  "scripts": {`,
						`    "test": "echo \"Error: no test specified\" && exit 1"`,
						`  },`,
						`  "keywords": [],`,
						`  "author": "",`,
						`  "license": "ISC",`,
						`  "dependencies": {`,
						`    "express": "^4.18.2"`,
						`  }`,
						`}`,
						``,
					},
				},
			},
			UngenCmd: `// UNGEN: inject file:package.json on ln 2 'replace "simple-nodejs" with kebabCase(var.appName)'`,
			ExpectedHistory: map[string][]int{
				"../examples/simple-nodejs/package.json": {2},
			},
			ExpectedContent: map[string][]string{
				"../examples/simple-nodejs/package.json": {
					`{`,
					`// UNGEN: replace "simple-nodejs" with kebabCase(var.appName)`,
					`  "name": "simple-nodejs",`,
					`  "version": "1.0.0",`,
					`  "description": "",`,
					`  "main": "index.js",`,
					`  "scripts": {`,
					`    "test": "echo \"Error: no test specified\" && exit 1"`,
					`  },`,
					`  "keywords": [],`,
					`  "author": "",`,
					`  "license": "ISC",`,
					`  "dependencies": {`,
					`    "express": "^4.18.2"`,
					`  }`,
					`}`,
					``,
				},
			},
		},
		{
			Context: InjectionContext{
				DotFilePath: "../examples/simple-nodejs/.ungen",
				InjectionHistory: map[string][]int{
					"../examples/simple-nodejs/package.json": {2},
				},
				InjectionContent: map[string][]string{
					"../examples/simple-nodejs/package.json": {
						`{`,
						`// UNGEN: replace "simple-nodejs" with kebabCase(var.appName)`,
						`  "name": "simple-nodejs",`,
						`  "version": "1.0.0",`,
						`  "description": "",`,
						`  "main": "index.js",`,
						`  "scripts": {`,
						`    "test": "echo \"Error: no test specified\" && exit 1"`,
						`  },`,
						`  "keywords": [],`,
						`  "author": "",`,
						`  "license": "ISC",`,
						`  "dependencies": {`,
						`    "express": "^4.18.2"`,
						`  }`,
						`}`,
						``,
					},
				},
			},
			UngenCmd: `// UNGEN: inject file:package.json on ln 5 'if var.useTypescript == "true" then replace ".js" with ".ts"'`,
			ExpectedHistory: map[string][]int{
				"../examples/simple-nodejs/package.json": {2, 5},
			},
			ExpectedContent: map[string][]string{
				"../examples/simple-nodejs/package.json": {
					`{`,
					`// UNGEN: replace "simple-nodejs" with kebabCase(var.appName)`,
					`  "name": "simple-nodejs",`,
					`  "version": "1.0.0",`,
					`  "description": "",`,
					`// UNGEN: if var.useTypescript == "true" then replace ".js" with ".ts"`,
					`  "main": "index.js",`,
					`  "scripts": {`,
					`    "test": "echo \"Error: no test specified\" && exit 1"`,
					`  },`,
					`  "keywords": [],`,
					`  "author": "",`,
					`  "license": "ISC",`,
					`  "dependencies": {`,
					`    "express": "^4.18.2"`,
					`  }`,
					`}`,
					``,
				},
			},
		},
	}

	for i, c := range testCases {
		cmd := c.UngenCmd
		p, _ := Parse(cmd)
		p.Inject(&c.Context)
		actualContent := c.Context.InjectionContent
		actualHistory := c.Context.InjectionHistory

		eqHistory := compareMaps(actualHistory, c.ExpectedHistory)
		if !eqHistory {
			t.Errorf("Test case %d failed for different history", i)
			break
		}
		fmt.Println("✔️ Same injection history for test case: ", i)
		eqContent := compareMaps(actualContent, c.ExpectedContent)
		if !eqContent {
			t.Errorf("Test case %d failed for different content", i)
			break
		}
		fmt.Println("✔️ Same injection content for test case: ", i)
		fmt.Println("✅", cmd)
	}
}

func compareMaps[T comparable](actual map[string][]T, expected map[string][]T) bool {
	equal := true
	for k, v1 := range actual {
		if v2, ok := expected[k]; ok {
			if !slicesEqual(v1, v2) {
				fmt.Printf("⏰ Different values for key %s:\n %v and\n %v\n", k, v1, v2)
				equal = false
			}
		} else {
			fmt.Printf("Key %s not found in expected\n", k)
			equal = false
		}
	}
	for k := range expected {
		if _, ok := actual[k]; !ok {
			fmt.Printf("Key %s not found in actual\n", k)
			equal = false
		}
	}

	return equal
}

func slicesEqual[T comparable](a, b []T) bool {

	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
