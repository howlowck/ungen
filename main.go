package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	inputDir := flag.String("i", "", "InputDirectory (Required)")
	flag.Parse()

	if *inputDir == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	filepath.Walk(*inputDir, func(path string, info os.FileInfo, err error) error {
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
			// Check if the line starts with "//" or "#"
			if strings.HasPrefix(line, "// UNGEN: ") || strings.HasPrefix(line, "# UNGEN: ") {
				fmt.Println(line)
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}

		return nil
	})
}
