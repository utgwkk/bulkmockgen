package generator

import (
	"context"
	"os"
	"os/exec"
)

type Runner interface {
	Run() error
}

var _ Runner = new(exec.Cmd)

func NewCommandRunner(ctx context.Context, cmdExecutable string, cmdArgs ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, cmdExecutable, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
