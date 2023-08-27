package internal

import (
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type LiteProgram struct {
	LiteCommands []*LiteCommand `@@*`
}

type LiteCommand struct {
	Header        *Header        `@@`
	LiteOperation *LiteOperation `@@`
}

type LiteOperation struct {
	Lines *NextLines `( @@`
	File  *ThisFile  `  | @@ )?`
	Text  *Text      `@@`
}

type NextLines struct {
	Lines *int `("next" @INT ( "lines" | "line" ))`
}

type ThisFile struct {
	ThisFile *bool `"for this file"`
}

type Text struct {
	Segments []*Segment `"@" @@* "@"`
}

type Segment struct {
	Value    *Value `( "${" @@ "}"`
	UserText string `  | @Char )`
}

func LiteParse(code string) (*LiteProgram, error) {
	program, err := liteParser.Parse("", strings.NewReader(code))
	if err != nil {
		return nil, err
	}
	return program, nil
}

var (
	// TODO: Have to use Stateful: https://github.com/alecthomas/participle/blob/master/_examples/stateful/main.go
	// liteLexer = lexer.MustSimple([]lexer.SimpleRule{
	// 	{"whitespace", `\s+`},
	// 	{"PAREN", `(\(|\))`},
	// 	{"COMMA", `,`},
	// 	{"CMDSTR", `'[^']*'`},
	// 	{"STR", `"[^"]*"`},
	// 	{"HEADER", `\S*\s?UNGEN:(\S+)? `},
	// 	{"KEYWORD", `(?i)\b(for|next|this)\b`},
	// 	{"INT", `\d+`},
	// 	{"UNIT", `(?i)\b(lines|line|file|folder)\b`},
	// 	{"STRFUNC", `(kebabCase|snakeCase|camelCase|upperCase|lowerCase|substitute|concat)\b`},
	// 	{"VAR", `var\.\w+`},
	// 	{"CLIPB", `cb\.\w+`},
	// 	{"RAW", `.+`},
	// })

	liteLexer = lexer.MustStateful(lexer.Rules{
		"Root": {
			{"STR", `"[^"]*"`, nil},
			{`Text`, "@", lexer.Push("Prompt")},
			{"HEADER", `\S*\s?UNGEN:(\S+)? `, nil},
			{"KEYWORD", `(?i)\b(for|next|this)\b`, nil},
			{"INT", `\d+`, nil},
			{"UNIT", `(?i)\b(lines|line|file|folder)\b`, nil},
			{"STRFUNC", `(kebabCase|snakeCase|camelCase|upperCase|lowerCase|substitute|concat)\b`, nil},
			{"VAR", `var\.\w+`, nil},
			{"CLIPB", `cb\.\w+`, nil},
		},
		"Text": {
			{`Char`, "\\$|[^$`\\\\]+", nil},
			{`Escaped`, `\\.`, nil},
			{`TextEnd`, "@", lexer.Pop()},
			{`Value`, `\${`, lexer.Push("Value")},
		},
		"Value": {
			{`Whitespace`, `\s+`, nil},
			{`Something`, `[^}\s]+`, nil},
			{`ValueEnd`, `}`, lexer.Pop()},
		},
	})

	liteParser = participle.MustBuild[LiteProgram](
		participle.Lexer(liteLexer),
		participle.CaseInsensitive("KEYWORD"),
		participle.Unquote("STR"),
		// participle.Unquote("CMDSTR"),
	)
)
