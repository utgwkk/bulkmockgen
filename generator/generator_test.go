package generator

import (
	"context"
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
				UseGoRun:    false,
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
				UseGoRun:    true,
				DryRun:      false,
				MockSetName: "MockInterfaces",
				RestArgs:    []string{"-package", "mock_multifile", "-destination", "mock_multifile/mock.go"},
			},
			wantCmdExecutable: "go",
			wantCmdArgs:       []string{"run", "go.uber.org/mock/mockgen", "-package", "mock_multifile", "-destination", "mock_multifile/mock.go", ".", "IFoo,IBar"},
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
