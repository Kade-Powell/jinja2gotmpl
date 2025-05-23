package transpiler

import (
	"strings"

	"github.com/Kade-Powell/jinja2gotmpl/internal/ast"
)

func ToGoTemplate(root *ast.RootNode) (string, error) {
	var sb strings.Builder
	for _, node := range root.Children {
		renderNode(&sb, node)
	}
	return sb.String(), nil
}

func renderNode(sb *strings.Builder, node ast.Node) {
	switch n := node.(type) {
	case *ast.TextNode:
		sb.WriteString(n.Content)
	case *ast.VariableNode:
		if len(n.Filters) == 0 {
			sb.WriteString("{{ .")
			sb.WriteString(n.Base)
			sb.WriteString(" }}")
		} else {
			sb.WriteString("{{ ")
			// Reverse the filter chain
			for i := len(n.Filters) - 1; i >= 0; i-- {
				sb.WriteString(n.Filters[i].Name)
				sb.WriteString(" ")
			}
			sb.WriteString(".")
			sb.WriteString(n.Base)
			sb.WriteString(" }}")
		}
	case *ast.IfNode:
		sb.WriteString("{{ if ")
		sb.WriteString(n.Condition)
		sb.WriteString(" }}")
		for _, child := range n.Body {
			renderNode(sb, child)
		}
		if len(n.ElseBody) > 0 {
			sb.WriteString("{{ else }}")
			for _, child := range n.ElseBody {
				renderNode(sb, child)
			}
		}
		sb.WriteString("{{ end }}")
	}
}
