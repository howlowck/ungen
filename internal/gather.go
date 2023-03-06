package internal

import "fmt"

func (p *Program) Gather(ctx *Context) {
	for _, c := range p.Commands {

		if c.Operation != nil {
			c.Operation.Gather(ctx)
		}

		if c.IfConditional != nil {
			conditionalValue := c.IfConditional.Condition.Evaluate(*ctx, []string{})
			if conditionalValue[0] == "true" {
				c.IfConditional.Then.Gather(ctx)
			}
		}
	}

}

func (v *Operation) Gather(ctx *Context) error {
	if v.Cut != nil {
		cbKey := v.Cut.To.Name
		if v.Cut.From.NextLines != nil {
			numOfLines := *v.Cut.From.NextLines
			startLIndex := ctx.ProgramLineNumber // the next line
			tempLines := ctx.Lines[startLIndex : startLIndex+numOfLines]
			ctx.Clipboard[cbKey] = tempLines
		}

		if v.Cut.From.LineNumRange != nil {
			count := v.Cut.From.LineNumRange.ToLn - v.Cut.From.LineNumRange.FromLn + 1

			if count < 1 {
				fmt.Println("You cannot have zero or negative lines in a line range")
			}
			startLIndex := v.Cut.From.LineNumRange.FromLn - 1
			numOfLines := count
			tempLines := ctx.Lines[startLIndex : startLIndex+numOfLines]
			ctx.Clipboard[cbKey] = tempLines
		}

		if v.Cut.From.LineNum != nil {
			startLIndex := *v.Cut.From.LineNum - 1
			tempLines := ctx.Lines[startLIndex:*v.Cut.From.LineNum]
			ctx.Clipboard[cbKey] = tempLines
		}
	}

	if v.Copy != nil {
		cbKey := v.Copy.To.Name

		if v.Copy.From.LineNumRange != nil {
			count := v.Copy.From.LineNumRange.ToLn - v.Copy.From.LineNumRange.FromLn + 1

			if count < 1 {
				fmt.Println("You cannot have zero or negative lines in a line range")
			}
			startLIndex := v.Copy.From.LineNumRange.FromLn - 1
			numOfLines := count
			tempLines := ctx.Lines[startLIndex : startLIndex+numOfLines]
			ctx.Clipboard[cbKey] = tempLines
		}
		if v.Copy.From.LineNum != nil {
			startLIndex := *v.Copy.From.LineNum - 1
			tempLines := ctx.Lines[startLIndex:1]
			ctx.Clipboard[cbKey] = tempLines
		}
		if v.Copy.From.NextLines != nil {
			numOfLines := *v.Copy.From.NextLines
			startLIndex := ctx.ProgramLineNumber // the next line
			tempLines := ctx.Lines[startLIndex : startLIndex+numOfLines]
			ctx.Clipboard[cbKey] = tempLines
		}
	}

	return nil
}
