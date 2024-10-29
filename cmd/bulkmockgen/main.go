package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"

	"slices"

	"github.com/utgwkk/bulkmockgen/generator"
)

var (
	useGoRun = flag.Bool("use_go_run", false, "Whether to use go run command to execute mockgen; defaults to false")
	dryRun   = flag.Bool("dry_run", false, "Print command to be executed and exit")
)

func main() {
	oldArgs := make([]string, len(os.Args))
	copy(oldArgs, os.Args)
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	mockSetName := flag.Arg(0)
	args := oldArgs
	restSeparatorIdx := slices.Index(oldArgs, "--")
	if restSeparatorIdx == -1 {
		log.Fatal("rest separator (--) is required")
	}

	g := &generator.Generator{
		UseGoRun:    *useGoRun,
		DryRun:      *dryRun,
		MockSetName: mockSetName,
		RestArgs:    args[restSeparatorIdx+1:],
	}
	rf := func(ctx context.Context, cmdExecutable string, cmdArgs ...string) generator.Runner {
		cmd := exec.CommandContext(ctx, cmdExecutable, cmdArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd
	}

	if err := g.Generate(ctx, rf); err != nil {
		log.Fatal(err)
	}
}
