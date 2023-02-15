package main

import (
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Program struct {
	Commands []*Command `@@*`
}

type Command struct {
	Header *Header `@@`

	IfConditional *If        `( @@`
	Operation     *Operation ` | @@ )`
}

type Header struct {
	HeaderText string `@HEADER`
}

type If struct {
	Condition *Expr      `"IF" @@`
	Then      *Operation `"THEN" @@`
	Else      *Operation `("ELSE" @@)?`
}

type Expr struct {
	Value *Value `@@`
}

type Operation struct {
	Replace *Replace `( @@`
	Delete  *Delete  ` | @@ )`
}

type Replace struct {
	From *Value `"replace" @@`
	To   *Value `"with" @@`
}

type Delete struct {
	NumOfLines int  `( "delete" @INT ( "lines" | "line" )`
	File       bool `  | "delete" @"file" `
	Directory  bool `  | "delete" @"folder" )`
}

type Value struct {
	// TODO: figure out how parse out the braces and value in string
	String   *string      `( @STR`
	Variable *string      `  | @VAR`
	StrFunc  *StrFunction `| @@ )`
}

type StrFunction struct {
	FunctionName string   `@STRFUNC`
	LeftParen    *string  `"("`
	Params       []*Value `@@? ("," @@)*`
	RightParen   *string  `")"`
}

func Parse(code string) (*Program, error) {
	program, err := basicParser.Parse("", strings.NewReader(code))
	if err != nil {
		return nil, err
	}
	return program, nil
}

var (
	basicLexer = lexer.MustSimple([]lexer.SimpleRule{
		{"whitespace", `\s+`},
		{"PAREN", `(\(|\))`},
		{"COMMA", `,`},
		{"STR", `'[^']*'|"[^"]*"`},
		{"STRFUNC", `(kebabCase|snakeCase|camelCase|upperCase|lowerCase|substitute)\b`},
		{"HEADER", `(\/\/|#) UNGEN:(v1)? `},
		{"INT", `\d+`},
		{"KEYWORD", `(?i)\b(if|then|else|replace|with|delete)\b`},
		{"UNIT", `(?i)\b(lines|line|file|folder)\b`},
		{"VAR", `var\.\w+`},
		{"EOL", `[\n\r]+`},
	})

	basicParser = participle.MustBuild[Program](
		participle.Lexer(basicLexer),
		participle.CaseInsensitive("KEYWORD"),
		participle.Unquote("STR"),
	)
)
