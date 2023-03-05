package main

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
				dotFilePath:      "examples/simple-nodejs/.ungen",
				injectionHistory: map[string][]int{},
				injectionContent: map[string][]string{},
			},
			UngenCmd: `// UNGEN: inject file:package.json on ln 2 'replace "simple-nodejs" with kebabCase(var.appName)'`,
			ExpectedHistory: map[string][]int{
				"examples/simple-nodejs/package.json": {2},
			},
			ExpectedContent: map[string][]string{
				"examples/simple-nodejs/package.json": {
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
				dotFilePath: "examples/simple-nodejs/.ungen",
				injectionHistory: map[string][]int{
					"examples/simple-nodejs/package.json": {2},
				},
				injectionContent: map[string][]string{
					"examples/simple-nodejs/package.json": {
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
				"examples/simple-nodejs/package.json": {2, 5},
			},
			ExpectedContent: map[string][]string{
				"examples/simple-nodejs/package.json": {
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
		actualContent := c.Context.injectionContent
		actualHistory := c.Context.injectionHistory

		eqHistory := compareMaps(actualHistory, c.ExpectedHistory)
		eqContent := compareMaps(actualContent, c.ExpectedContent)
		if !eqHistory {
			fmt.Println(actualHistory, c.ExpectedHistory)
			t.Errorf("Test case %d failed for different history", i)
		}
		if !eqContent {
			fmt.Println(actualContent, c.ExpectedContent)
			t.Errorf("Test case %d failed for different content", i)
		}
		fmt.Println("✅", cmd)
	}
}

func compareMaps[T comparable](actual map[string][]T, expected map[string][]T) bool {
	equal := true
	for k, v1 := range actual {
		if v2, ok := expected[k]; ok {
			if !slicesEqual(v1, v2) {
				fmt.Printf("Different values for key %s:\n %v and\n %v\n", k, v1, v2)
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
