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
				// TODO: recursively handle Value
				replaceTo = c.Operation.Replace.To.Evaluate(ctx, []string{})
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

func (v *Value) Evaluate(ctx EvalContext, inputs []string) string {
	if v.String != nil {
		return *v.String
	}
	if v.Variable != nil {
		varTemp := *v.Variable
		varName := varTemp[4:] // take away 'var.'
		return ctx.vars[varName]
	}
	// TODO: add string transformation functions
	return ""
}
