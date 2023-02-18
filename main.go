package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

type varMap map[string]string

// Implement Set method for varMap
func (m *varMap) Set(s string) error {
	kv := strings.Split(s, "=")
	if len(kv) == 2 {
		(*m)[kv[0]] = kv[1]
	}
	// TODO: maybe need to unset a value of there is 1 element
	return nil
}

// Implement String method for kvMap
func (m *varMap) String() string {
	return fmt.Sprint(*m)
}

func main() {
	vars := make(varMap)

	inputDir := flag.String("i", "", "InputDirectory (Required)")
	outputDir := flag.String("o", "", "OutputDirectory (Required)")
	keepLine := flag.Bool("keep", false, "Keep the UNGEN line")
	flag.Var(&vars, "var", "Set Variables (ex. -var foo=bar -var baz=qux)")

	flag.Parse()

	if *inputDir == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *outputDir == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	r, _ := regexp.Compile(`\s?[\/]?[\/|#] UNGEN: (.*)\s?$`)

	tempDir, err := ioutil.TempDir(os.TempDir(), "ungen-")

	if err != nil {
		panic(err)
	}

	ignoreList := getIgnorePatterns(filepath.Join(*inputDir, ".gitignore"))

	// 1. Copy the directory into a staging directory
	err = copyDir(*inputDir, tempDir, ignoreList)

	if err != nil {
		panic(err)
	}

	fmt.Println("copied to Dir:", tempDir)
	filepath.Walk(tempDir, func(path string, info os.FileInfo, e error) error {
		// Skip directories (since they will be scanned recursively)
		if info.IsDir() {
			return nil
		}

		// Read the file
		body, err := os.ReadFile(path)
		if err != nil {
			// Handle error
			log.Fatalf("unable to read file: %v", err)
		}

		lines := strings.Split(string(body), "\n")

		// 2. Process in the staging directory
		for i, v := range lines {
			// fmt.Println(i, v)
			if r.MatchString(v) {
				context := EvalContext{
					lines:             lines,
					vars:              vars,
					path:              path,
					keepLine:          *keepLine,
					programLineNumber: i + 1,
				}
				fmt.Println(v)
				program, _ := Parse(v)
				patches := program.Evaluate(context)
				for _, patch := range patches {
					if patch.Content != nil {
						lines = patch.Content.Apply(lines)
					}
				}
			}
		}

		// Overwrite the file with new content
		os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0777)
		return nil
	})

	// 3. Copy the staging directory to the output directory
	copyDir(tempDir, *outputDir, ignoreList)
}

func getIgnorePatterns(path string) []string {
	// TODO: need to get the .gitignore file from every level
	// get .gitignore file
	gitIgnore, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer gitIgnore.Close()

	// Create a new scanner to read the contents of the file
	scanner := bufio.NewScanner(gitIgnore)

	// Create a slice to store the ignore patterns
	ignore := make([]string, 0)

	// Iterate through the lines of the file
	for scanner.Scan() {
		// Add each line as an ignore pattern
		ignore = append(ignore, strings.TrimSpace(scanner.Text()))
	}

	// check for any errors while scanning
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return ignore
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func copyDir(src string, dst string, ignoreList []string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		for _, pattern := range ignoreList {
			normalizedPattern := strings.TrimLeft(strings.Trim(pattern, "*"), "/")
			// base := filepath.Base(path)
			match, _ := doublestar.PathMatch("**/"+normalizedPattern+"**", path)

			if match {
				// fmt.Println("pattern:", pattern, "path:", path, "match:", match)
				return nil
			}
		}

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		} else {
			return copyFile(path, dstPath)
		}
	})
}
