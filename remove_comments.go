package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		return
	}

	filename := os.Args[1]

	// Create the file set
	fset := token.NewFileSet()

	// Parse the file
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return
	}

	// Clear comments from the AST
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.File:
			x.Comments = nil
		case *ast.GenDecl:
			x.Doc = nil
			x.Specs = clearCommentFromSpecs(x.Specs)
		case *ast.FuncDecl:
			x.Doc = nil
		case *ast.Field:
			x.Doc = nil
			x.Comment = nil
		case *ast.ValueSpec:
			x.Doc = nil
			x.Comment = nil
		case *ast.TypeSpec:
			x.Doc = nil
			x.Comment = nil
		case *ast.ImportSpec:
			x.Doc = nil
			x.Comment = nil
		case *ast.FieldList:
			if x != nil {
				for _, field := range x.List {
					field.Doc = nil
					field.Comment = nil
				}
			}
		}
		return true
	})

	// Write the modified file back
	out, err := os.Create(filename)
	if err != nil {
		return
	}
	defer out.Close()

	err = printer.Fprint(out, fset, node)
	if err != nil {
		return
	}
}

func clearCommentFromSpecs(specs []ast.Spec) []ast.Spec {
	for _, spec := range specs {
		switch x := spec.(type) {
		case *ast.ValueSpec:
			x.Doc = nil
			x.Comment = nil
		case *ast.TypeSpec:
			x.Doc = nil
			x.Comment = nil
		}
	}
	return specs
}