package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// OsExitAnalyzer анализатор проверяющий наличие метода os.Exit в пакете main.
var OsExitAnalyzer = &analysis.Analyzer{
	Name: "exitAnalyzer",
	Doc:  "Don't allow os.Exit in main package",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.File:
				if x.Name.Name != "main" {
					return false
				}
			case *ast.SelectorExpr:
				if x.Sel.Name == "Exit" {
					pass.Reportf(x.Pos(), "expression has os.Exit call in main package")

				}
			}
			return true
		})
	}
	return nil, nil
}
