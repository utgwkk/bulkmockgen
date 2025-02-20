package generator

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

const (
	mockgenPackage = "go.uber.org/mock/mockgen"
)

type ExecMode int

const (
	ExecModeDirect ExecMode = iota + 1

	ExecModeGoRun
)

func (e ExecMode) Command() (string, []string) {
	switch e {
	case ExecModeDirect:
		return "mockgen", []string{}
	case ExecModeGoRun:
		return "go", []string{"run", mockgenPackage}
	default:
		panic("unsupported exec mode")
	}
}

type runnerFactory func(ctx context.Context, cmdExecutable string, cmdArgs ...string) Runner

type Generator struct {
	ExecMode    ExecMode
	DryRun      bool
	MockSetName string
	RestArgs    []string
}

func (g *Generator) Generate(ctx context.Context, rf runnerFactory) error {
	mockSet, externalPkg, err := g.findMockSet(".")
	if err != nil {
		return err
	}

	cmdExecutable, cmdArgs := g.ExecMode.Command()
	cmdArgs = append(cmdArgs, g.RestArgs...)
	if externalPkg != "" {
		cmdArgs = append(cmdArgs, externalPkg)
	} else {
		cmdArgs = append(cmdArgs, ".")
	}
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

func (g *Generator) findMockSet(sourceDir string) ([]string, string, error) {
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return g.findMockSetFromPackage(sourceDir)
	}
	return g.findMockSetFromDirectory(sourceDir)
}

func (g *Generator) findMockSetFromPackage(packagePath string) ([]string, string, error) {
	pp, err := packages.Load(&packages.Config{}, packagePath)
	if err != nil {
		return nil, "", err
	}
	pkg := pp[0]
	return g.findMockSetFromDirectory(pkg.Module.Dir)
}

func (g *Generator) findMockSetFromDirectory(sourceDir string) ([]string, string, error) {
	fset := token.NewFileSet()
	parsed, err := parser.ParseDir(fset, sourceDir, nil, 0)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse: %w", err)
	}
	for _, p := range parsed {
		for _, f := range p.Files {
			mockSet, externalPkg, err := lookupFromScope(fset, f.Scope, g.MockSetName)
			if err != nil {
				if errors.Is(err, errMockSetNotFound) {
					continue
				}
				return nil, "", err
			}
			if externalPkg != "" {
				for _, imp := range f.Imports {
					impPath, _ := strconv.Unquote(imp.Path.Value)
					if imp.Name != nil && externalPkg == imp.Name.Name {
						externalPkg = impPath
						break
					}

					splittedImpPath := strings.Split(impPath, "/")
					if externalPkg == splittedImpPath[len(splittedImpPath)-1] {
						externalPkg = impPath
						break
					}
				}
			}
			return mockSet, externalPkg, nil
		}
	}

	return nil, "", fmt.Errorf("mock set %s not found", g.MockSetName)
}

var errMockSetNotFound = errors.New("mock set notfound")

func lookupFromScope(fset *token.FileSet, s *ast.Scope, mockSetName string) ([]string, string, error) {
	obj := s.Lookup(mockSetName)
	if obj == nil {
		return nil, "", errMockSetNotFound
	}
	if obj.Kind != ast.Var {
		return nil, "", errMockSetNotFound
	}
	vs, ok := obj.Decl.(*ast.ValueSpec)
	if !ok {
		return nil, "", errMockSetNotFound
	}

	var externalPackageName string
	var mockSet []string
	for _, v := range vs.Values {
		cmp, ok := v.(*ast.CompositeLit)
		if !ok {
			return nil, "", newInvalidMockSetError("mock set is not a slice", fset, v)
		}
		for _, v := range cmp.Elts {
			funCall, ok := v.(*ast.CallExpr)
			if !ok {
				return nil, "", newInvalidMockSetError("mock set contains not a function call", fset, v)
			}
			// new(...)
			newIdent, ok := funCall.Fun.(*ast.Ident)
			if !ok {
				return nil, "", newInvalidMockSetError("mock set contains not a new(...) function call", fset, v)
			}
			if newIdent.Name != "new" {
				return nil, "", newInvalidMockSetError("mock set contains not a new(...) function call", fset, v)
			}
			if len(funCall.Args) != 1 {
				return nil, "", newInvalidMockSetError("mock set contains not a new(...) function call", fset, v)
			}
			arg := funCall.Args[0]
			switch arg := arg.(type) {
			case *ast.Ident:
				mockSet = append(mockSet, arg.Name)
			case *ast.SelectorExpr:
				x := arg.X
				ident, ok := x.(*ast.Ident)
				if !ok {
					return nil, "", newInvalidMockSetError("not an external package selector", fset, arg)
				}
				if externalPackageName != "" && externalPackageName != ident.Name {
					return nil, "", newInvalidMockSetError("mixing external package interfaces is not allowed", fset, arg)
				}
				externalPackageName = ident.Name
				mockSet = append(mockSet, arg.Sel.Name)
			default:
				return nil, "", newInvalidMockSetError("not an interface is passed to new(...) function call", fset, arg)
			}
		}
	}
	return mockSet, externalPackageName, nil
}

func newInvalidMockSetError(msg string, fset *token.FileSet, e ast.Expr) error {
	pos := fset.Position(e.Pos())
	return fmt.Errorf("%s at file %s, line %d, column %d", msg, pos.Filename, pos.Line, pos.Column)
}
