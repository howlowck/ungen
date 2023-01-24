package main

import (
	"fmt"
	"regexp"
	"strings"
)

type PatchType int

const (
	PatchReplace = iota
	PatchDelete
	PatchInsert
)

func (t PatchType) String() string {
	return [...]string{"Replace", "Delete", "Insert"}[t]
}

type Patch struct {
	PatchType     PatchType
	OldLineNumber int
	OldLineCount  int
	NewContent    string
}

func Eval(fileText string, vars map[string]string, p *Program, programLineNumber int) []Patch {
	// get all args
	lines := regexp.MustCompile("\r?\n").Split(fileText, -1)

	// get old Line Number and Range
	oldLineNumber := programLineNumber + 1 // the next line
	oldLineCount := 1                      // default to 1

	// create patch from program

	oldContent := strings.Join(lines[oldLineNumber-1:oldLineNumber+oldLineCount-1], "\n")
	fmt.Println(oldContent)
	patch := Patch{
		PatchType:     PatchDelete,
		OldLineNumber: oldLineNumber,
		OldLineCount:  oldLineCount,
		NewContent:    "",
	}
	return []Patch{patch}
}

func (p *Program) Evaluate(lines []string, vars map[string]string, programLineNumber int) []Patch {

	// get old Line Number and Range

	var result Patch
	for _, c := range p.Commands {

		if c.Operation != nil {
			if c.Operation.Replace != nil {
				replaceFrom := c.Operation.Replace.From.String

				var replaceTo *string
				if c.Operation.Replace.To.String != nil {
					replaceTo = c.Operation.Replace.To.String
				} else {
					// TODO evaluate variable
					varTemp := *c.Operation.Replace.To.Variable
					varName := varTemp[4:]
					varValue := vars[varName]
					replaceTo = &varValue
				}
				oldLineNumber := programLineNumber + 1 // the next line
				oldLineCount := 1                      // default to 1

				oldContent := strings.Join(lines[oldLineNumber-1:oldLineNumber+oldLineCount-1], "\n")
				re := regexp.MustCompile(*replaceFrom)
				newContent := re.ReplaceAllString(oldContent, *replaceTo)
				result = Patch{
					PatchType:     PatchReplace,
					OldLineNumber: oldLineNumber,
					OldLineCount:  oldLineCount,
					NewContent:    newContent,
				}
				fmt.Print(result)
			}
		}
	}
	return []Patch{result}
}
