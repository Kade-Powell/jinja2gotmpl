package j2g

import (
	"github.com/Kade-Powell/jinja2gotmpl/internal/lexer"
	"github.com/Kade-Powell/jinja2gotmpl/internal/parser"
	"github.com/Kade-Powell/jinja2gotmpl/internal/transpiler"
)

func Transpile(input string) (string, error) {
	lex := lexer.New(input)
	tokens := []lexer.Token{}
	for {
		tok := lex.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == lexer.TokenEOF {
			break
		}
	}

	astParser := parser.New(tokens)
	root, err := astParser.Parse()
	if err != nil {
		return "", err
	}

	goTemplate, err := transpiler.ToGoTemplate(root)
	if err != nil {
		return "", err
	}
	return goTemplate, nil
}
