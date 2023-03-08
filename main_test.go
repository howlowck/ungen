package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestMain(t *testing.T) {
	os.Args = []string{
		"ungen",
		"-i", "examples/simple-nodejs",
		"-o", "test",
		"-var", "includeExtraFeature=true",
		"-var", "appName=Haos Awesome App",
		"-var", "theme=dark",
		"-var", "useTypescript=false",
	}
	fmt.Printf("\nos args = %v\n", os.Args)
	main()

	fmt.Println(`
=================================================
Begin comparing generated files to expected files
=================================================
	`)

	fmt.Println("🔍 Comparing File Tree...")

	isDirsSame := compFileTree("test", "examples/expected-simple-nodejs-nokeep")
	if !isDirsSame {
		t.Errorf("generated files are not the same as the expected ones")
		fmt.Println("🚨 Directories have different file trees")
		os.RemoveAll("test")
		return
	}

	fmt.Println("✅ File Trees are the same for both directories")

	fmt.Println("-----------------------------------------------")

	fmt.Println("🔍 Comparing File Content...")
	if !compFiles("test", "examples/expected-simple-nodejs-nokeep") {
		t.Errorf("generated content are not the same as the expected ones")
		fmt.Println("🚨 File contents are not the same for both directories")
		os.RemoveAll("test")
		return
	}

	fmt.Println("✅ File contents are the same for both directories")
	os.RemoveAll("test")
}

func compContents(actualContent, expectedContent, path string) bool {
	re := regexp.MustCompile(`\r?\n`)
	actualLines := re.Split(actualContent, -1)
	expectedLines := re.Split(expectedContent, -1)
	result := true
	for i, line := range actualLines {
		if line != expectedLines[i] {
			fmt.Printf("✖️ In %s, Ln %d does not match: \n %s\n %s\n\n", path, i, line, expectedLines[i])
			result = false
		}
	}
	return result
}

func compFiles(actualPath, expectedPath string) bool {
	files1, err := os.ReadDir(actualPath)
	if err != nil {
		log.Fatal(err)
	}

	files2, err := os.ReadDir(expectedPath)
	if err != nil {
		log.Fatal(err)
	}

	result := true
	for i := range files1 {
		if files1[i].IsDir() && !compFiles(filepath.Join(actualPath, files1[i].Name()), filepath.Join(expectedPath, files2[i].Name())) {
			result = false
		}
		if !files1[i].IsDir() {
			actualFilePath := filepath.Join(actualPath, files1[i].Name())
			expectedFilePath := filepath.Join(expectedPath, files2[i].Name())
			actualContent, _ := os.ReadFile(actualFilePath)
			expectedContent, _ := os.ReadFile(expectedFilePath)
			isContentSame := compContents(string(actualContent), string(expectedContent), filepath.Join(actualPath, files1[i].Name())[5:])

			if !isContentSame {
				result = false
			}
			// fmt.Println("✔️ " + actualFilePath + " and " + expectedFilePath + " have the same content")
		}
	}
	// if result {
	// fmt.Println("✔️✔️ " + actualPath + " and " + expectedPath + " have the same contents")
	// }
	return result
}

func compFileTree(actualPath, expectedPath string) bool {
	files1, err := os.ReadDir(actualPath)
	if err != nil {
		log.Fatal(err)
	}

	files2, err := os.ReadDir(expectedPath)
	if err != nil {
		log.Fatal(err)
	}

	actualNum := len(files1)
	expectedNum := len(files2)
	if actualNum != expectedNum {
		fmt.Printf("✖️ %s and %s have different number of files: %d and %d respectively.\n", actualPath, expectedPath, actualNum, expectedNum)
		return false
	}

	for i := range files1 {
		if files1[i].Name() != files2[i].Name() {
			fmt.Println("✖️ " + actualPath + " and " + expectedPath + " have different file names: " + files1[i].Name() + " and " + files2[i].Name())
			return false
		}
		if files1[i].IsDir() && !compFileTree(filepath.Join(actualPath, files1[i].Name()), filepath.Join(expectedPath, files2[i].Name())) {
			return false
		}
	}

	// fmt.Println("✔️✔️ " + actualPath + " and " + expectedPath + " have the same file names")
	return true
}
