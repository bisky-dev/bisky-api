package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	fset := token.NewFileSet()
	var violations []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() {
			switch d.Name() {
			case ".git", "vendor", "tmp", "docs", "tools":
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}
		base := filepath.Base(path)
		if base == "types.go" || strings.HasSuffix(base, "_test.go") {
			return nil
		}

		normalizedPath := filepath.ToSlash(path)
		if strings.HasPrefix(normalizedPath, "internal/db/sqlc/") || strings.Contains(normalizedPath, "/internal/db/sqlc/") {
			return nil
		}

		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		if isGenerated(file) {
			return nil
		}

		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				pos := fset.Position(typeSpec.Pos())
				violations = append(violations, fmt.Sprintf("%s:%d type %q must be in types.go", pos.Filename, pos.Line, typeSpec.Name.Name))
			}
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "typefilelint failed: %v\n", err)
		os.Exit(1)
	}

	if len(violations) > 0 {
		fmt.Fprintln(os.Stderr, "typefilelint found type declarations outside types.go:")
		for _, v := range violations {
			fmt.Fprintln(os.Stderr, " -", v)
		}
		os.Exit(1)
	}
}

func isGenerated(file *ast.File) bool {
	if file.Doc == nil {
		return false
	}
	for _, c := range file.Doc.List {
		text := strings.ToLower(c.Text)
		if strings.Contains(text, "code generated") && strings.Contains(text, "do not edit") {
			return true
		}
	}
	return false
}
