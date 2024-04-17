package osexit

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for use to os.exit",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				if x.Name.Name == "main" {
					nodes := hasOsExit(x)
					for _, a := range nodes {
						pass.Reportf(a.Pos(), "it is recommended to avoid calling the os.Exit function from the main function")
					}
				}
			}
			return true
		})
	}

	return nil, nil
}

func hasOsExit(x ast.Node) []*ast.CallExpr {
	result := make([]*ast.CallExpr, 0)
	ast.Inspect(x, func(node ast.Node) bool {
		switch y := node.(type) {
		case *ast.CallExpr:
			if s, ok := y.Fun.(*ast.SelectorExpr); ok {
				if i, ok := s.X.(*ast.Ident); ok {
					if i.Name == "os" && s.Sel.Name == "Exit" {
						result = append(result, y)
					}
				}
			}
		}
		return true
	})
	return result
}
