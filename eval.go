package main

import (
	"regexp"
	"strings"
)

// TODO: figure out how to do file or directory deletions
func (p *Program) Evaluate(lines []string, vars map[string]string, programLineNumber int) []Patch {

	var result []Patch

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
				patch := Patch{
					Content: &ContentPatch{
						PatchType:     PatchReplace,
						OldLineNumber: oldLineNumber,
						OldLineCount:  oldLineCount,
						NewContent:    newContent,
					}}

				result = append(result, patch)

				// fmt.Print(result)
			}
			if c.Operation.Delete != nil {
				oldLineNumber := programLineNumber + 1 // the next line
				oldLineCount := c.Operation.Delete.NumOfLines
				contentPatch := ContentPatch{
					PatchType:     PatchDelete,
					OldLineNumber: oldLineNumber,
					OldLineCount:  oldLineCount,
					NewContent:    []string{},
				}
				patch := Patch{
					Content: &contentPatch,
				}
				// fmt.Println(patch)
				result = append(result, patch)
			}
		}
	}
	return result
}
