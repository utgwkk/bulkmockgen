package generator

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

const (
	mockgenPackage = "go.uber.org/mock/mockgen"
)

type runnerFactory func(ctx context.Context, cmdExecutable string, cmdArgs ...string) Runner

type Generator struct {
	UseGoRun    bool
	DryRun      bool
	MockSetName string
	RestArgs    []string
}

func (g *Generator) Generate(ctx context.Context, rf runnerFactory) error {
	mockSet, err := g.findMockSet(".")
	if err != nil {
		return err
	}

	var cmdExecutable string
	var cmdArgs []string
	if g.UseGoRun {
		cmdExecutable = "go"
		cmdArgs = append(cmdArgs, "run", mockgenPackage)
	} else {
		cmdExecutable = "mockgen"
	}
	cmdArgs = append(cmdArgs, g.RestArgs...)
	cmdArgs = append(cmdArgs, ".")
	cmdArgs = append(cmdArgs, strings.Join(mockSet, ","))
	if g.DryRun {
		fmt.Printf("%s %s", cmdExecutable, strings.Join(cmdArgs, " "))
		return nil
	}
	runner := rf(ctx, cmdExecutable, cmdArgs...)
	if err := runner.Run(); err != nil {
		return err
	}

	return nil
}

func (g *Generator) findMockSet(sourceDir string) ([]string, error) {
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return g.findMockSetFromPackage(sourceDir)
	}
	return g.findMockSetFromDirectory(sourceDir)
}

func (g *Generator) findMockSetFromPackage(packagePath string) ([]string, error) {
	pp, err := packages.Load(&packages.Config{}, packagePath)
	if err != nil {
		return nil, err
	}
	pkg := pp[0]
	return g.findMockSetFromDirectory(pkg.Module.Dir)
}

func (g *Generator) findMockSetFromDirectory(sourceDir string) ([]string, error) {
	fset := token.NewFileSet()
	parsed, err := parser.ParseDir(fset, sourceDir, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}
	for _, p := range parsed {
		for _, f := range p.Files {
			mockSet, err := lookupFromScope(f.Scope, g.MockSetName)
			if err != nil {
				if errors.Is(err, errMockSetNotFound) {
					continue
				}
				return nil, err
			}
			return mockSet, nil
		}
	}

	return nil, fmt.Errorf("mock set %s not found", g.MockSetName)
}

var errMockSetNotFound = errors.New("mock set notfound")

func lookupFromScope(s *ast.Scope, mockSetName string) ([]string, error) {
	obj := s.Lookup(mockSetName)
	if obj == nil {
		return nil, errMockSetNotFound
	}
	if obj.Kind != ast.Var {
		return nil, errMockSetNotFound
	}
	vs, ok := obj.Decl.(*ast.ValueSpec)
	if !ok {
		return nil, errMockSetNotFound
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
