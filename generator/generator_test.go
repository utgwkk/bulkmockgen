package generator

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockedRunner struct {
	cmdExecutable string
	cmdArgs       []string
}

var _ Runner = new(mockedRunner)

func (*mockedRunner) Run() error {
	return nil
}

func TestGenerate(t *testing.T) {
	testcases := []struct {
		name              string
		baseDir           string
		g                 *Generator
		wantCmdExecutable string
		wantCmdArgs       []string
	}{
		{
			name:    "basic",
			baseDir: "fixtures/multifile",
			g: &Generator{
				ExecMode:    ExecModeDirect,
				DryRun:      false,
				MockSetName: "MockInterfaces",
				RestArgs:    []string{"-package", "mock_multifile", "-destination", "mock_multifile/mock.go"},
			},
			wantCmdExecutable: "mockgen",
			wantCmdArgs:       []string{"-package", "mock_multifile", "-destination", "mock_multifile/mock.go", ".", "IFoo,IBar"},
		},
		{
			name:    "with go run",
			baseDir: "fixtures/multifile",
			g: &Generator{
				ExecMode:    ExecModeGoRun,
				DryRun:      false,
				MockSetName: "MockInterfaces",
				RestArgs:    []string{"-package", "mock_multifile", "-destination", "mock_multifile/mock.go"},
			},
			wantCmdExecutable: "go",
			wantCmdArgs:       []string{"run", "go.uber.org/mock/mockgen", "-package", "mock_multifile", "-destination", "mock_multifile/mock.go", ".", "IFoo,IBar"},
		},
		{
			name:    "with go tool",
			baseDir: "fixtures/multifile",
			g: &Generator{
				ExecMode:    ExecModeGoTool,
				DryRun:      false,
				MockSetName: "MockInterfaces",
				RestArgs:    []string{"-package", "mock_multifile", "-destination", "mock_multifile/mock.go"},
			},
			wantCmdExecutable: "go",
			wantCmdArgs:       []string{"tool", "go.uber.org/mock/mockgen", "-package", "mock_multifile", "-destination", "mock_multifile/mock.go", ".", "IFoo,IBar"},
		},
		{
			name:    "with external package",
			baseDir: "fixtures/external",
			g: &Generator{
				ExecMode:    ExecModeDirect,
				DryRun:      false,
				MockSetName: "MockInterfaces",
				RestArgs:    []string{"-package", "mock_sql", "-destination", "mock_sql/mock.go"},
			},
			wantCmdExecutable: "mockgen",
			wantCmdArgs:       []string{"-package", "mock_sql", "-destination", "mock_sql/mock.go", "database/sql/driver", "Conn,Driver"},
		},
	}
	for _, tc := range testcases {
		tc := tc
		origDir, err := os.Getwd()
		require.NoError(t, err)
		t.Run(tc.name, func(t *testing.T) {
			err := os.Chdir(tc.baseDir)
			defer os.Chdir(origDir)
			require.NoError(t, err)
			var r *mockedRunner
			newMockedRunner := func(ctx context.Context, cmdExecutable string, cmdArgs ...string) Runner {
				r = &mockedRunner{
					cmdExecutable: cmdExecutable,
					cmdArgs:       cmdArgs,
				}
				return r
			}

			err = tc.g.Generate(context.Background(), newMockedRunner)
			require.NoError(t, err)
			require.NotNil(t, r)
			assert.Equal(t, tc.wantCmdExecutable, r.cmdExecutable)
			assert.Equal(t, tc.wantCmdArgs, r.cmdArgs)
		})
	}
}

func prepareScope(t *testing.T, file string) (*token.FileSet, *ast.Scope) {
	t.Helper()

	fset := token.NewFileSet()
	parsed, err := parser.ParseFile(fset, file, nil, 0)
	require.NoError(t, err)

	return fset, parsed.Scope
}

func TestLookupFromScope(t *testing.T) {
	testcases := []struct {
		name             string
		inputFile        string
		inputMockSetName string
		wantErrMsg       string
	}{
		{
			name:             "not found",
			inputFile:        "./fixtures/ng/no_mock_set/prog.go",
			inputMockSetName: "NotFound",
			wantErrMsg:       "mock set notfound",
		},
		{
			name:             "invalid expression",
			inputFile:        "./fixtures/ng/invalid_expr/prog.go",
			inputMockSetName: "Iset",
			wantErrMsg:       "mock set contains not a function call at file ./fixtures/ng/invalid_expr/prog.go, line 4, column 2",
		},
		{
			name:             "not a new function call",
			inputFile:        "./fixtures/ng/not_a_new_function/prog.go",
			inputMockSetName: "Iset",
			wantErrMsg:       "mock set contains not a new(...) function call at file ./fixtures/ng/not_a_new_function/prog.go, line 4, column 2",
		},
		{
			name:             "new takes not an interface",
			inputFile:        "./fixtures/ng/not_an_interface/prog.go",
			inputMockSetName: "Iset",
			wantErrMsg:       "not an interface is passed to new(...) function call at file ./fixtures/ng/not_an_interface/prog.go, line 4, column 6",
		},
		{
			name:             "not a slice",
			inputFile:        "./fixtures/ng/not_a_slice/prog.go",
			inputMockSetName: "Iset",
			wantErrMsg:       "mock set is not a slice at file ./fixtures/ng/not_a_slice/prog.go, line 3, column 12",
		},
		{
			name:             "mixed external packages",
			inputFile:        "./fixtures/ng/mixed_external/prog.go",
			inputMockSetName: "Iset",
			wantErrMsg:       "mixing external package interfaces is not allowed at file ./fixtures/ng/mixed_external/prog.go, line 10, column 6",
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			fset, s := prepareScope(t, tc.inputFile)
			_, _, err := lookupFromScope(fset, s, tc.inputMockSetName)
			require.Error(t, err)
			assert.Equal(t, tc.wantErrMsg, err.Error())
		})
	}
}
