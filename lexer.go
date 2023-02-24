package main

import (
	"regexp"
	"strconv"
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
	Condition *Value     `"if" @@`
	Then      *Operation `"then" @@`
	Else      *Operation `("else" @@)?`
}

type Operation struct {
	Replace *Replace `( @@`
	Copy    *Copy    ` | @@`
	Cut     *Cut     ` | @@`
	Insert  *Insert  ` | @@`
	Delete  *Delete  ` | @@ )`
}

type Insert struct {
	Value *Value `"insert" @@`
}

type Replace struct {
	From *Value `"replace" @@`
	To   *Value `"with" @@`
}

type Delete struct {
	NumOfLines *int  `( "delete" @INT ( "lines" | "line" )`
	File       *bool `  | "delete" @"file" `
	Directory  *bool `  | "delete" @"folder" )`
}

type Copy struct {
	From *ContentLines `"copy" @@`
	To   *ClipBoard    `"to" @@`
}

type Cut struct {
	From *ContentLines `"cut" @@`
	To   *ClipBoard    `"to" @@`
}

type ContentLines struct {
	NextLines    *int          `( "next" @INT ( "lines" | "line" )`
	LineNum      *int          `  | "ln" @INT `
	LineNumRange *LineNumRange `  | "ln" @@ )`
}

type LineNumRange struct {
	FromLn int `@@`
	ToLn   int `@@`
}

func (v *LineNumRange) Parse(lex *lexer.PeekingLexer) error {
	regex, _ := regexp.Compile(`^\d+\-\d+`)
	tok := lex.Peek()
	if !regex.MatchString(tok.Value) {
		return participle.NextMatch
	}
	lex.Next()
	numbers := strings.Split(tok.Value, "-")
	fromNum, _ := strconv.Atoi(numbers[0])
	toNum, _ := strconv.Atoi(numbers[1])
	*v = LineNumRange{
		FromLn: fromNum,
		ToLn:   toNum,
	}
	return nil
}

type Value struct {
	// TODO: figure out how parse out the braces and value in string
	String    *string      `( @STR`
	StrFunc   *StrFunction `  | @@ `
	Variable  *Variable    `  | @@`
	ClipBoard *ClipBoard   `  | @@ )`
}

type Variable struct {
	Name string `@@`
}

func (v *Variable) Parse(lex *lexer.PeekingLexer) error {
	tok := lex.Peek()
	if !strings.HasPrefix(tok.Value, "var.") {
		return participle.NextMatch
	}
	lex.Next()
	*v = Variable{
		Name: tok.Value[4:],
	}
	return nil
}

type ClipBoard struct {
	Name string `@@`
}

func (v *ClipBoard) Parse(lex *lexer.PeekingLexer) error {
	tok := lex.Peek()
	if !strings.HasPrefix(tok.Value, "cb.") {
		return participle.NextMatch
	}
	lex.Next()
	*v = ClipBoard{
		Name: tok.Value[3:],
	}
	return nil
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
		{"NUMRANGE", `\d+\-\d+`},
		{"STRFUNC", `(kebabCase|snakeCase|camelCase|upperCase|lowerCase|substitute|concat)\b`},
		{"HEADER", `(\/\/|#) UNGEN:(v1)? `},
		{"INT", `\d+`},
		{"KEYWORD", `(?i)\b(if|then|else|replace|with|delete|copy|cut|to|insert|next|ln)\b`},
		{"UNIT", `(?i)\b(lines|line|file|folder)\b`},
		{"VAR", `var\.\w+`},
		{"CLIPB", `cb\.\w+`},
		{"EOL", `[\n\r]+`},
	})

	basicParser = participle.MustBuild[Program](
		participle.Lexer(basicLexer),
		participle.CaseInsensitive("KEYWORD"),
		participle.Unquote("STR"),
	)
)
