package main

import (
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

// TODO: figure out how to do file or directory deletions
func (p *Program) Evaluate(lines []string, vars map[string]string, programLineNumber int) []ContentPatch {

	var result []ContentPatch

	for _, c := range p.Commands {

		if c.Operation != nil {
			if c.Operation.Replace != nil {
				replaceFrom := c.Operation.Replace.From.String

				var replaceTo string
				if c.Operation.Replace.To.String != nil {
					replaceTo = *(c.Operation.Replace.To.String)
				} else {
					varTemp := *c.Operation.Replace.To.Variable
					varName := varTemp[4:] // take away 'var.'
					replaceTo = vars[varName]
				}
				oldLineNumber := programLineNumber + 1 // the next line
				oldLineCount := 1                      // default to 1

				oldContent := strings.Join(lines[oldLineNumber-1:oldLineNumber+oldLineCount-1], "\n")
				re := regexp.MustCompile(*replaceFrom)
				newContent := strings.Split(re.ReplaceAllString(oldContent, replaceTo), "\n")
				patch := ContentPatch{
					PatchType:     PatchReplace,
					OldLineNumber: oldLineNumber,
					OldLineCount:  oldLineCount,
					NewContent:    newContent,
				}

				result = append(result, patch)

				// fmt.Print(result)
			}
			if c.Operation.Delete != nil {
				oldLineNumber := programLineNumber + 1 // the next line
				oldLineCount := c.Operation.Delete.NumOfLines
				patch := ContentPatch{
					PatchType:     PatchDelete,
					OldLineNumber: oldLineNumber,
					OldLineCount:  oldLineCount,
					NewContent:    []string{},
				}
				// fmt.Println(patch)
				result = append(result, patch)
			}
		}
	}
	return result
}
