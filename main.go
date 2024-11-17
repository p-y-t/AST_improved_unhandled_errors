package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"
)

func visitNodes(fset *token.FileSet, pkg *types.Package, info *types.Info, n ast.Node) {
	ast.Inspect(n, func(n ast.Node) bool {
		switch stmt := n.(type) {
		case *ast.AssignStmt:
			for idx, lhs := range stmt.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok && ident.Name == "_" {
					for _, rhs := range stmt.Rhs {
						if callExpr, ok := rhs.(*ast.CallExpr); ok {
							printFunctionCallWithReturnType(fset, info, callExpr, idx)
						}
					}
				}
			}
		}
		return true
	})
}

func printFunctionCallWithReturnType(fset *token.FileSet, info *types.Info, callExpr *ast.CallExpr, paramIndex int) {
	funcName := "<unknown>"
	pos := fset.Position(callExpr.Pos())

	switch fun := callExpr.Fun.(type) {
	case *ast.Ident:
		funcName = fun.Name
	case *ast.SelectorExpr:
		if pkg, ok := fun.X.(*ast.Ident); ok {
			funcName = fmt.Sprintf("%s.%s", pkg.Name, fun.Sel.Name)
		}
	}

	if typ := info.TypeOf(callExpr); typ != nil {
		if funcSig, ok := typ.(*types.Tuple); ok {
			if paramIndex < funcSig.Len() {
				returnType := funcSig.At(paramIndex).Type()
				fmt.Printf("Line %d: Assignment to '_': Function call: %s, Param index: %d, Return type: %s\n",
					pos.Line, funcName, paramIndex+1, returnType)
			} else {
				fmt.Printf("Line %d: Assignment to '_': Function call: %s, Param index: %d, Return type: <out of range>\n",
					pos.Line, funcName, paramIndex+1)
			}
		} else {
			if paramIndex == 0 {
				fmt.Printf("Line %d: Assignment to '_': Function call: %s, Param index: 1, Return type: %s\n",
					pos.Line, funcName, typ.String())
			} else {
				fmt.Printf("Line %d: Assignment to '_': Function call: %s, Param index: %d, Return type: <invalid>\n",
					pos.Line, funcName, paramIndex+1)
			}
		}
	} else {
		fmt.Printf("Line %d: Assignment to '_': Function call: %s, Param index: %d, Return type: <unknown>\n",
			pos.Line, funcName, paramIndex+1)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <directory>")
		return
	}

	dir := os.Args[1]
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			fmt.Printf("Parsing file: %s\n", path)

			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
			if err != nil {
				fmt.Printf("Failed to parse %s: %v\n", path, err)
				return nil
			}

			conf := types.Config{
				Importer: importer.For("source", nil),
				Error: func(err error) {
					fmt.Printf("Type-checking error: %v\n", err)
				},
			}

			typeInfo := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
			}

			_, err = conf.Check(dir, fset, []*ast.File{node}, typeInfo)
			if err != nil {
				fmt.Printf("Type-checking failed for %s: %v\n", path, err)
				return nil
			}

			visitNodes(fset, nil, typeInfo, node)
			fmt.Println()
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the directory: %v\n", err)
	}
}
