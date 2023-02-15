package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/bmatcuk/doublestar/v4"
)

func main() {
	inputDir := flag.String("i", "", "InputDirectory (Required)")
	outputDir := flag.String("o", "", "OutputDirectory (Required)")
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

		// Open the file
		file, err := os.Open(path)
		// fmt.Println("path:", path)
		if err != nil {
			return err
		}
		defer file.Close()

		// 2. Process in the staging directory
		// Scan the file line by line
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			if r.MatchString(line) {
				// Process Line
				// Extract the line after the prefix
				cmd := r.FindStringSubmatch(line)
				fmt.Println(cmd[1])
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}

		// 3. Copy the staging directory to the output directory
		err = copyDir(tempDir, *outputDir, ignoreList)

		return nil
	})
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
		ignore = append(ignore, scanner.Text())
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
			// base := filepath.Base(path)
			match, _ := doublestar.PathMatch(pattern, path)

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
