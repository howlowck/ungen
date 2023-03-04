package main

import (
	"fmt"
	"path"
	"reflect"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

type Context struct {
	lines             []string
	path              string
	vars              map[string]string
	clipboard         map[string][]string
	keepLine          bool
	programLineNumber int
}

func ProcessLineNumber(patch Patch, keepLine bool) Patch {
	result := patch

	if !keepLine && patch.Content.PatchType == PatchInsert {
		result.Content.OldLineNumber = patch.Content.OldLineNumber - 1
		result.Content.OldLineCount = patch.Content.OldLineCount
	} else if !keepLine {
		result.Content.OldLineNumber = patch.Content.OldLineNumber - 1
		result.Content.OldLineCount = patch.Content.OldLineCount + 1
	}

	return result
}

// TODO: figure out how to do file or directory deletions
func (p *Program) Evaluate(ctx Context) []Patch {

	var result []Patch

	for _, c := range p.Commands {

		if c.Operation != nil {
			result = append(result, c.Operation.Evaluate(ctx)...)
		}

		if c.IfConditional != nil {
			conditionalValue := c.IfConditional.Condition.Evaluate(ctx, []string{})
			if conditionalValue[0] == "true" {
				result = append(result, c.IfConditional.Then.Evaluate(ctx)...)
			} else {
				if c.IfConditional.Else != nil {
					result = append(result, c.IfConditional.Else.Evaluate(ctx)...)
				} else {
					oldLineNumber := ctx.programLineNumber + 1 // the next line
					contentPatch := ContentPatch{
						PatchType:     PatchDelete,
						OldLineNumber: oldLineNumber,
						OldLineCount:  0,
						NewContent:    []string{},
					}
					patch := Patch{
						Content: &contentPatch,
					}
					processed := ProcessLineNumber(patch, ctx.keepLine)
					result = append(result, processed)
				}
			}
		}
	}

	return result
}

func (c *Conditional) Evaluate(ctx Context, vars []string) []string {
	left := c.Left.Evaluate(ctx, vars)
	if c.Right != nil {
		right := c.Right.Evaluate(ctx, vars)
		isEqual := reflect.DeepEqual(left, right)
		if isEqual && *c.Op == "==" {
			return []string{"true"}
		} else if !isEqual && *c.Op == "!=" {
			return []string{"true"}
		}
		return []string{"false"}
	}

	return left
}

func (v *Operation) Evaluate(ctx Context) []Patch {
	if v.Replace != nil {
		replaceFrom := v.Replace.From.String

		replaceTo := v.Replace.To.Evaluate(ctx, []string{})
		oldLineNumber := ctx.programLineNumber + 1 // the next line
		oldLineCount := 1                          // default to 1

		oldContent := strings.Join(ctx.lines[oldLineNumber-1:oldLineNumber+oldLineCount-1], "\n")
		re := regexp.MustCompile(*replaceFrom)
		newContent := strings.Split(re.ReplaceAllString(oldContent, replaceTo[0]), "\n")
		patch := Patch{
			Content: &ContentPatch{
				PatchType:     PatchReplace,
				OldLineNumber: oldLineNumber,
				OldLineCount:  oldLineCount,
				NewContent:    newContent,
			}}

		return []Patch{ProcessLineNumber(patch, ctx.keepLine)}
	}

	if v.Copy != nil {
		if v.Copy.From.LineNumRange != nil {
			count := v.Copy.From.LineNumRange.ToLn - v.Copy.From.LineNumRange.FromLn + 1

			if count < 1 {
				fmt.Println("You cannot have zero or negative lines in a line range")
			}
			if !ctx.keepLine {
				return []Patch{{
					Content: &ContentPatch{
						PatchType:     PatchDelete,
						OldLineNumber: ctx.programLineNumber,
						OldLineCount:  1,
						NewContent:    []string{},
					}},
				}
			}
			return []Patch{}
		}
		if v.Copy.From.LineNum != nil {
			if !ctx.keepLine {
				return []Patch{{
					Content: &ContentPatch{
						PatchType:     PatchDelete,
						OldLineNumber: ctx.programLineNumber,
						OldLineCount:  1,
						NewContent:    []string{},
					}},
				}
			}
			return []Patch{}
		}
		if v.Copy.From.NextLines != nil {
			patch := Patch{
				Content: &ContentPatch{
					PatchType:     PatchDelete,
					OldLineNumber: ctx.programLineNumber + 1,
					OldLineCount:  0,
					NewContent:    []string{},
				}}

			return []Patch{ProcessLineNumber(patch, ctx.keepLine)}
		}
	}

	if v.Cut != nil {
		if v.Cut.From.NextLines != nil {
			patch := Patch{
				Content: &ContentPatch{
					PatchType:     PatchDelete,
					OldLineNumber: ctx.programLineNumber + 1,
					OldLineCount:  0,
					NewContent:    []string{},
				}}

			return []Patch{ProcessLineNumber(patch, ctx.keepLine)}
		}

		if v.Cut.From.LineNumRange != nil {
			count := v.Cut.From.LineNumRange.ToLn - v.Cut.From.LineNumRange.FromLn + 1

			if count < 1 {
				fmt.Println("You cannot have zero or negative lines in a line range")
			}
			patch := Patch{
				Content: &ContentPatch{
					PatchType:     PatchDelete,
					OldLineNumber: v.Cut.From.LineNumRange.FromLn,
					OldLineCount:  count,
					NewContent:    []string{},
				}}
			result := []Patch{patch}
			if !ctx.keepLine {
				remoteCommandLine := Patch{
					Content: &ContentPatch{
						PatchType:     PatchDelete,
						OldLineNumber: ctx.programLineNumber,
						OldLineCount:  1,
						NewContent:    []string{},
					}}
				result = append(result, remoteCommandLine)
			}
			return result
		}

		if v.Cut.From.LineNum != nil {
			patch := Patch{
				Content: &ContentPatch{
					PatchType:     PatchDelete,
					OldLineNumber: *v.Cut.From.LineNum,
					OldLineCount:  1,
					NewContent:    []string{},
				}}
			result := []Patch{patch}
			if !ctx.keepLine {
				remoteCommandLine := Patch{
					Content: &ContentPatch{
						PatchType:     PatchDelete,
						OldLineNumber: ctx.programLineNumber,
						OldLineCount:  1,
						NewContent:    []string{},
					}}
				result = append(result, remoteCommandLine)
			}
			return result
		}
	}

	if v.Insert != nil {
		value := v.Insert.Value.Evaluate(ctx, []string{})
		patch := Patch{
			Content: &ContentPatch{
				PatchType:     PatchReplace,
				OldLineNumber: ctx.programLineNumber + 1,
				OldLineCount:  0,
				NewContent:    value,
			}}
		return []Patch{ProcessLineNumber(patch, ctx.keepLine)}
	}

	if v.Delete.File != nil {
		patch := Patch{
			File: &FilePatch{
				FileOp:     FileDelete,
				TargetPath: ctx.path,
			},
		}
		return []Patch{patch}
	}

	if v.Delete.Directory != nil {
		dir, _ := path.Split(ctx.path)
		patch := Patch{
			File: &FilePatch{
				FileOp:     DirectoryDelete,
				TargetPath: dir,
			},
		}
		return []Patch{patch}
	}

	// if not replace, then it's delete
	oldLineNumber := ctx.programLineNumber + 1 // the next line
	oldLineCount := v.Delete.NumOfLines
	contentPatch := ContentPatch{
		PatchType:     PatchDelete,
		OldLineNumber: oldLineNumber,
		OldLineCount:  *oldLineCount,
		NewContent:    []string{},
	}
	patch := Patch{
		Content: &contentPatch,
	}

	return []Patch{ProcessLineNumber(patch, ctx.keepLine)}

}

func (v *Value) Evaluate(ctx Context, inputs []string) []string {
	if v.String != nil {
		return []string{*v.String}
	}

	if v.Variable != nil {
		varTemp := *v.Variable
		varName := varTemp.Name
		strValue, ok := ctx.vars[varName]
		if ok {
			return []string{strValue}
		} else {
			fmt.Println("warning! Variable " + varName + " does not exist")
			return []string{}
		}
	}

	if v.ClipBoard != nil {
		cbTemp := *v.ClipBoard
		cbName := cbTemp.Name
		strValue, ok := ctx.clipboard[cbName]
		if ok {
			return strValue
		} else {
			fmt.Println("warning! Clipboard " + cbName + " does not exist")
			return []string{}
		}
	}

	if v.StrFunc != nil {
		params := []string{}

		for _, pv := range v.StrFunc.Params {
			value := pv.Evaluate(ctx, []string{})
			params = append(params, value[0])
		}

		// TODO what to put in the slice?
		value := params[0]
		if v.StrFunc.FunctionName == "kebabCase" {
			return []string{strcase.ToKebab(value)}
		}
		if v.StrFunc.FunctionName == "snakeCase" {
			return []string{strcase.ToSnake(value)}
		}
		if v.StrFunc.FunctionName == "camelCase" {
			return []string{strcase.ToLowerCamel(value)}
		}
		if v.StrFunc.FunctionName == "upperCamelCase" {
			return []string{strcase.ToCamel(value)}
		}
		if v.StrFunc.FunctionName == "upperCase" {
			return []string{strings.ToUpper(value)}
		}
		if v.StrFunc.FunctionName == "lowerCase" {
			return []string{strings.ToLower(value)}
		}
		if v.StrFunc.FunctionName == "substitute" {
			return []string{strings.ReplaceAll(value, params[1], params[2])}
		}
		if v.StrFunc.FunctionName == "concat" {
			var result string
			for _, curr := range params {
				result += curr
			}
			return []string{result}
		}
	}
	return []string{}
}
