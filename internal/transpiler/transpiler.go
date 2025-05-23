package transpiler

import (
	"strings"

	"github.com/Kade-Powell/jinja2gotmpl/internal/ast"
)

func ToGoTemplate(root *ast.RootNode) (string, error) {
	var sb strings.Builder
	localVars := map[string]bool{}
	for _, node := range root.Children {
		renderNode(&sb, node, "", localVars)
	}
	return sb.String(), nil
}

func renderNode(sb *strings.Builder, node ast.Node, loopVar string, localVars map[string]bool) {
	switch n := node.(type) {
	case *ast.TextNode:
		sb.WriteString(n.Content)

	case *ast.VariableNode:
		sb.WriteString("{{ ")

		// Apply filters
		for i := len(n.Filters) - 1; i >= 0; i-- {
			filter := n.Filters[i]
			sb.WriteString(filter.Name)
			sb.WriteString(" ")
			for _, arg := range filter.Args {
				sb.WriteString(arg)
				sb.WriteString(" ")
			}
		}

		// Determine variable scope
		switch {
		case loopVar != "" && strings.HasPrefix(n.Base, loopVar+"."):
			sb.WriteString("$")
			sb.WriteString(loopVar)
			sb.WriteString(".")
			sb.WriteString(strings.TrimPrefix(n.Base, loopVar+"."))
		case strings.HasPrefix(n.Base, "loop."):
			sb.WriteString("$loop.")
			sb.WriteString(strings.TrimPrefix(n.Base, "loop."))
		case loopVar != "" && strings.HasPrefix(n.Base, loopVar+"."):
			sb.WriteString("$")
			sb.WriteString(loopVar)
			sb.WriteString(".")
			sb.WriteString(strings.TrimPrefix(n.Base, loopVar+"."))
		case localVars[n.Base]:
			sb.WriteString("$")
			sb.WriteString(n.Base)
		default:
			sb.WriteString(".")
			sb.WriteString(n.Base)
		}

		sb.WriteString(" }}")

	case *ast.IfNode:
		sb.WriteString("{{ if ")
		sb.WriteString(n.Condition)
		sb.WriteString(" }}")
		for _, child := range n.Body {
			renderNode(sb, child, loopVar, localVars)
		}
		if len(n.ElseBody) > 0 {
			sb.WriteString("{{ else }}")
			for _, child := range n.ElseBody {
				renderNode(sb, child, loopVar, localVars)
			}
		}
		sb.WriteString("{{ end }}")

	case *ast.ForNode:
		indexVar := "$i"
		itemVar := "$" + n.Item
		listExpr := "." + n.List

		// range with index
		sb.WriteString("{{ range ")
		sb.WriteString(indexVar)
		sb.WriteString(", ")
		sb.WriteString(itemVar)
		sb.WriteString(" := ")
		sb.WriteString(listExpr)
		sb.WriteString(" }}\n")

		// Inject $loop object (if needed)
		sb.WriteString("{{ $loop := dict \"index\" (add ")
		sb.WriteString(indexVar)
		sb.WriteString(" 1) \"last\" (eq (add ")
		sb.WriteString(indexVar)
		sb.WriteString(" 1) (len ")
		sb.WriteString(listExpr)
		sb.WriteString(")) }}\n")

		// Transpile body
		for _, child := range n.Body {
			renderNode(sb, child, n.Item, map[string]bool{"loop": true})
		}

		sb.WriteString("{{ end }}")

	case *ast.SetNode:
		localVars[n.Name] = true
		sb.WriteString("{{ $")
		sb.WriteString(n.Name)
		sb.WriteString(" := ")
		sb.WriteString(n.Value)
		sb.WriteString(" }}")
	}

}
