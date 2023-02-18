package main

import (
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

type EvalContext struct {
	lines             []string
	path              string
	vars              map[string]string
	programLineNumber int
}

func ProcessLineNumber(patch Patch, deleteLineNumber bool) Patch {
	result := patch

	if deleteLineNumber == true {
		result.Content.OldLineNumber = patch.Content.OldLineNumber - 1
		result.Content.OldLineCount = patch.Content.OldLineCount + 1
	}
	return result
}

// TODO: figure out how to do file or directory deletions
func (p *Program) Evaluate(ctx EvalContext) []Patch {

	var result []Patch

	for _, c := range p.Commands {

		if c.Operation != nil {
			result = append(result, c.Operation.Evaluate(ctx))
		}

		if c.IfConditional != nil {
			conditionalValue := c.IfConditional.Condition.Evaluate(ctx, []string{})
			if conditionalValue == "true" {
				result = append(result, c.IfConditional.Then.Evaluate(ctx))
			} else {
				result = append(result, c.IfConditional.Else.Evaluate(ctx))
			}
		}
	}

	return result
}

func (v *Operation) Evaluate(ctx EvalContext) Patch {
	if v.Replace != nil {
		replaceFrom := v.Replace.From.String

		replaceTo := v.Replace.To.Evaluate(ctx, []string{})
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

		// TODO: parse keep line or not
		return ProcessLineNumber(patch, false)
	}

	// if not replace, then it's delete
	oldLineNumber := ctx.programLineNumber + 1 // the next line
	oldLineCount := v.Delete.NumOfLines
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
	return ProcessLineNumber(patch, false)

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
	if v.StrFunc != nil {
		params := []string{}

		for _, pv := range v.StrFunc.Params {
			value := pv.Evaluate(ctx, []string{})
			params = append(params, value)
		}

		// TODO what to put in the slice?
		value := params[0]
		if v.StrFunc.FunctionName == "kebabCase" {
			return strcase.ToKebab(value)
		}
		if v.StrFunc.FunctionName == "snakeCase" {
			return strcase.ToSnake(value)
		}
		if v.StrFunc.FunctionName == "camelCase" {
			return strcase.ToLowerCamel(value)
		}
		if v.StrFunc.FunctionName == "upperCamelCase" {
			return strcase.ToCamel(value)
		}
		if v.StrFunc.FunctionName == "upperCase" {
			return strings.ToUpper(value)
		}
		if v.StrFunc.FunctionName == "lowerCase" {
			return strings.ToLower(value)
		}
		if v.StrFunc.FunctionName == "substitute" {
			return strings.ReplaceAll(value, params[1], params[2])
		}
		if v.StrFunc.FunctionName == "concat" {
			var result string
			for _, curr := range params {
				result += curr
			}
			return result
		}
	}
	return ""
}
