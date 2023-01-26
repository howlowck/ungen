package main

import (
	"regexp"
	"strings"
)

type EvalContext struct {
	lines             []string
	path              string
	vars              map[string]string
	programLineNumber int
}

// TODO: figure out how to do file or directory deletions
func (p *Program) Evaluate(ctx EvalContext) []Patch {

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
					replaceTo = ctx.vars[varName]
				}
				oldLineNumber := ctx.programLineNumber + 1 // the next line
				oldLineCount := 1                          // default to 1

				oldContent := strings.Join(ctx.lines[oldLineNumber-1:oldLineNumber+oldLineCount-1], "\n")
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
				oldLineNumber := ctx.programLineNumber + 1 // the next line
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
				// fmt.Println(contentPatch)
				result = append(result, patch)
			}
		}
	}
	return result
}
