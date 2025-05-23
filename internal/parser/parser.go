package parser

import (
	"fmt"
	"strings"

	"github.com/Kade-Powell/jinja2gotmpl/internal/ast"
	"github.com/Kade-Powell/jinja2gotmpl/internal/lexer"
)

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) current() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{Type: lexer.TokenEOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) next() lexer.Token {
	tok := p.current()
	p.pos++
	return tok
}

func (p *Parser) Parse() (*ast.RootNode, error) {
	root := &ast.RootNode{}
	for {
		tok := p.current()
		if tok.Type == lexer.TokenEOF {
			break
		}
		switch tok.Type {
		case lexer.TokenText:
			root.Children = append(root.Children, &ast.TextNode{Content: tok.Literal})
			p.next()
		case lexer.TokenVariable:
			varNode := parseVariable(tok.Literal)
			root.Children = append(root.Children, varNode)
			p.next()
		case lexer.TokenBlock:
			ifNode, err := p.parseIf()
			if err != nil {
				return nil, err
			}
			root.Children = append(root.Children, ifNode)
		default:
			p.next() // skip unknown
		}
	}
	return root, nil
}

func (p *Parser) parseIf() (*ast.IfNode, error) {
	tok := p.next() // consume {% if ... %}
	condition := strings.TrimSpace(strings.TrimPrefix(tok.Literal, "if"))
	ifNode := &ast.IfNode{
		Condition: condition,
		Body:      []ast.Node{},
		ElseBody:  []ast.Node{},
	}

	currentBody := &ifNode.Body

	for {
		tok := p.current()
		if tok.Type == lexer.TokenEOF {
			return nil, fmt.Errorf("unexpected EOF in if block")
		}

		if tok.Type == lexer.TokenBlock {
			switch strings.TrimSpace(tok.Literal) {
			case "else":
				p.next() // consume {% else %}
				currentBody = &ifNode.ElseBody
				continue
			case "endif":
				p.next() // consume {% endif %}
				return ifNode, nil
			default:
				// Future support for nested if, for, etc. — for now, skip
				p.next()
				continue
			}
		}

		switch tok.Type {
		case lexer.TokenText:
			*currentBody = append(*currentBody, &ast.TextNode{Content: tok.Literal})
		case lexer.TokenVariable:
			*currentBody = append(*currentBody, parseVariable(tok.Literal))
		}
		p.next()
	}
}

func parseVariable(input string) *ast.VariableNode {
	parts := strings.Split(input, "|")
	base := strings.TrimSpace(parts[0])
	var filters []ast.FilterCall

	for _, part := range parts[1:] {
		name := strings.TrimSpace(part)
		if name != "" {
			filters = append(filters, ast.FilterCall{Name: name})
		}
	}

	return &ast.VariableNode{
		Base:    base,
		Filters: filters,
	}
}
