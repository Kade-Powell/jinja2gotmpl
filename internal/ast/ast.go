package ast

import "strings"

type NodeType string

const (
	NodeText     NodeType = "Text"
	NodeVariable NodeType = "Variable"
	NodeIf       NodeType = "If"
	NodeFor      NodeType = "For"
	NodeBlock    NodeType = "Block"
	NodeComment  NodeType = "Comment"
	NodeRoot     NodeType = "Root"
)

type Node interface {
	Type() NodeType
	String() string
}

// TextNode represents plain text
type TextNode struct {
	Content string
}

func (n *TextNode) Type() NodeType { return NodeText }
func (n *TextNode) String() string { return n.Content }

type FilterCall struct {
	Name string
	Args []string
}

// VariableNode represents {{ variable | filters }}
type VariableNode struct {
	Base    string       // e.g. user.name
	Filters []FilterCall // e.g. upper, join, etc.
}

func (n *VariableNode) Type() NodeType { return NodeVariable }

func (n *VariableNode) String() string {
	var sb strings.Builder
	sb.WriteString("{{ ")
	sb.WriteString(n.Base)
	for _, f := range n.Filters {
		sb.WriteString(" | " + f.Name)
	}
	sb.WriteString(" }}")
	return sb.String()
}

// IfNode represents {% if ... %} ... {% else %} ... {% endif %}
type IfNode struct {
	Condition string
	Body      []Node
	ElseBody  []Node
}

func (n *IfNode) Type() NodeType { return NodeIf }
func (n *IfNode) String() string { return "{% if " + n.Condition + " %} ... {% endif %}" }

// RootNode is the top-level node holding the entire template
type RootNode struct {
	Children []Node
}

func (n *RootNode) Type() NodeType { return NodeRoot }
func (n *RootNode) String() string { return "Root" }
