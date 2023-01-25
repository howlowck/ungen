package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	inputDir := flag.String("i", "", "InputDirectory (Required)")
	flag.Parse()

	if *inputDir == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	r, _ := regexp.Compile(`\s?[\/]?[\/|#] UNGEN: (.*)\s?$`)

	filepath.Walk(*inputDir, func(path string, info os.FileInfo, e error) error {
		// Skip directories (since they will be scanned recursively)
		if info.IsDir() {
			return nil
		}

		// Open the file
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

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

		return nil
	})
}
