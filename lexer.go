package main

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Program struct {
	Commands []*Command `@@*`
}

type Command struct {
	Header *Header `@@`

	Operation *Operation ` @@`
}

type Header struct {
	HeaderText string `@HEADER`
}

// type If struct {
// 	Pos lexer.Position

// 	Condition *Value     `"IF" @@`
// 	Operation *Operation `"THEN" @@`
// }

type Operation struct {
	Replace *Replace `@@`
	// Delete  *Delete  `@@`
}

type Replace struct {
	From *Value `"replace" @@`
	To   *Value `"with" @@`
}

type Delete struct {
	NumOfLines *int `"delete" @Int "lines"`
}

type Value struct {
	String   string `( @STR`
	Variable string `  | @VAR )`
}

func Parse(code string) (*Program, error) {
	program, err := basicParser.Parse("", strings.NewReader(code))
	if err != nil {
		return nil, err
	}
	program.init()
	return program, nil
}

func (p *Program) init() {
	for _, cmd := range p.Commands {
		fmt.Println(cmd)
	}
}

var (
	basicLexer = lexer.MustSimple([]lexer.SimpleRule{
		{"whitespace", `\s+`},

		{"STR", `'[^']*'|"[^"]*"`},
		{"HEADER", `(\/\/|#) UNGEN:(v1)? `},
		{"INT", `\d+`},
		{"KEYWORD", `(?i)\b(if|replace|with|delete|lines)\b`},
		{"VAR", `var\.\w+`},
		{"EOL", `[\n\r]+`},
	})

	basicParser = participle.MustBuild[Program](
		participle.Lexer(basicLexer),
		participle.CaseInsensitive("KEYWORD"),
		participle.Unquote("STR"),
		// participle.UseLookahead(2),
	)
)
