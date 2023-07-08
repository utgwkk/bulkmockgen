package generator

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	mockgenPackage = "go.uber.org/mock/mockgen"
)

type Generator struct {
	UseGoRun      bool
	PackageName   string
	SourcePackage string
	MockSetName   string
	FileOut       string

	writer io.Writer // for testing
}

func (g *Generator) Generate() error {
	mockSet, err := g.findMockSet(g.SourcePackage)
	if err != nil {
		return err
	}
	abspath, err := filepath.Abs(g.SourcePackage)
	if err != nil {
		return err
	}
	splittedPath := strings.Split(filepath.ToSlash(abspath), "/")
	sourcePkg := splittedPath[len(splittedPath)-1]

	destination := filepath.Join(g.PackageName, g.PackageName+".go")
	mockgen := "mockgen"
	if g.UseGoRun {
		mockgen = "go run " + mockgenPackage
	}
	var f io.Writer
	if g.writer != nil {
		f = g.writer
	} else if g.FileOut == "" {
		f = os.Stdout
	} else {
		file, err := os.Create(g.FileOut)
		if err != nil {
			return err
		}
		defer file.Close()
		f = file
	}
	writeOutTo(f, sourcePkg, &dumpCommentOption{
		mockgenCmd:     mockgen,
		destination:    destination,
		packageName:    g.PackageName,
		sourcePackage:  ".",
		interfaceNames: mockSet,
	})

	return nil
}

func (g *Generator) findMockSet(sourceDir string) ([]string, error) {
	fset := token.NewFileSet()
	parsed, err := parser.ParseDir(fset, sourceDir, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}
	for _, p := range parsed {
		for _, f := range p.Files {
			x := f.Scope
			obj := x.Lookup(g.MockSetName)
			if obj == nil {
				continue
			}
			if obj.Kind != ast.Var {
				continue
			}
			vs, ok := obj.Decl.(*ast.ValueSpec)
			if !ok {
				continue
			}

			var mockSet []string
			for _, v := range vs.Values {
				cmp, ok := v.(*ast.CompositeLit)
				if !ok {
					return nil, errors.New("not a interface slice")
				}
				for _, v := range cmp.Elts {
					funCall, ok := v.(*ast.CallExpr)
					if !ok {
						return nil, errors.New("invalid mock set (contains non-new(...) expression)")
					}
					// new(...)
					newIdent, ok := funCall.Fun.(*ast.Ident)
					if !ok {
						return nil, errors.New("invalid mock set (contains non-new(...) expression)")
					}
					if newIdent.Name != "new" {
						return nil, errors.New("invalid mock set (contains non-new(...) expression)")
					}
					if len(funCall.Args) != 1 {
						return nil, errors.New("invalid mock set (contains non-new(...) expression)")
					}
					arg := funCall.Args[0]
					argIdent, ok := arg.(*ast.Ident)
					if !ok {
						return nil, errors.New("invalid mock set (contains non-new(...) expression)")
					}
					if len(funCall.Args) != 1 {
						return nil, errors.New("invalid mock set (contains non-new(...) expression)")
					}
					mockSet = append(mockSet, argIdent.Name)
				}
			}
			return mockSet, nil
		}
	}

	return nil, fmt.Errorf("mock set %s not found", g.MockSetName)
}
