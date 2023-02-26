package main

import (
	"archive/zip"
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
	zipOutput := flag.Bool("zip", false, "Zip the output directory into a file")

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

	r, _ := regexp.Compile(`^\s?[\/]?[\/|#] UNGEN: (.*)\s?$`)

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

	fmt.Println("✅ Copied to Temp Dir:", tempDir)

	// 2. Gather Stage (gather all the text into clipboard)
	clipboard := make(map[string][]string)
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

		for i, v := range lines {
			if r.MatchString(v) {
				context := Context{
					lines:             lines,
					vars:              vars,
					path:              path,
					keepLine:          *keepLine,
					clipboard:         clipboard,
					programLineNumber: i + 1,
				}
				program, error := Parse(v)
				if error != nil {
					fmt.Println("Error parsing line: " + v)
					fmt.Println(error)
					os.Exit(1)
				}
				program.Gather(&context)
			}
		}

		return nil
	})

	fmt.Println("✅ Completed Gather Stage")

	for i, v := range clipboard {
		fmt.Println("Clipboard: " + i)
		for _, vv := range v {
			fmt.Println("├─ " + vv)
		}
	}

	// 3. Eval and Patch Stage
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
		fileOps := []Patch{}
		fmt.Println("Processing file for Eval and Patch: " + strings.Replace(path, tempDir+"/", "", 1))
		for i, v := range lines {
			if r.MatchString(v) {
				fmt.Println("├─ Ungen Found: " + strings.TrimSpace(v))
				context := Context{
					lines:             lines,
					vars:              vars,
					path:              path,
					keepLine:          *keepLine,
					clipboard:         clipboard,
					programLineNumber: i + 1,
				}
				program, _ := Parse(v)
				patches := program.Evaluate(context)
				for _, patch := range patches {
					if patch.Content != nil {
						lines = patch.Content.Apply(lines)
					}
					if patch.File != nil {
						fileOps = append(fileOps, patch)
					}
				}
			}
		}

		// Overwrite the file with new content
		os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0777)

		// TODO: it's doing extra work.. need to exit fast later

		for _, p := range fileOps {
			if p.File != nil {
				if p.File.FileOp == FileDelete {
					os.Remove(p.File.TargetPath)
				}
				if p.File.FileOp == DirectoryDelete {
					os.RemoveAll(p.File.TargetPath)
				}
			}
		}

		fmt.Println("(Applying filesystem changes)")
		fmt.Println("=========== Done ============")
		fmt.Println("")
		return nil
	})

	// 3. Copy the staging directory to the output directory
	if *zipOutput {
		zipDir(tempDir, *outputDir)
		fmt.Println("🎉 Created zip file: " + strings.TrimRight(*outputDir, "/") + ".zip")
	} else {
		// TODO: have a -clean flag to delete the output directory first
		copyDir(tempDir, *outputDir, ignoreList)
		fmt.Println("🎉 Created a directory in " + strings.TrimRight(*outputDir, "/"))
	}
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
			normalizedPattern := strings.TrimLeft(pattern, "/")

			// append "**" only if pattern ends in /
			if strings.HasSuffix(pattern, "/") {
				normalizedPattern += "**"
			}

			match, _ := doublestar.PathMatch("**/"+normalizedPattern, path)

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

func zipDir(src string, dst string) {
	outFile, err := os.Create(dst + ".zip")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)

	filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		inFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer inFile.Close()

		// Add file to zip archive using relative path as name
		withSlashPath := strings.TrimRight(src, "/") + "/"
		cleanPath := strings.Replace(path, withSlashPath, "", 1)
		writer, err := zipWriter.Create(cleanPath)
		if err != nil {
			return err
		}

		// Write file content to zip archive using io.Copy method
		io.Copy(writer, inFile)

		return nil
	})
	zipWriter.Close()
}
