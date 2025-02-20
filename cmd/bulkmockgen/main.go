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

	flagExecMode = flag.String("exec_mode", "direct", "How to execute mockgen. Supported values are: direct, go_run, go_tool")

	dryRun = flag.Bool("dry_run", false, "Print command to be executed and exit")
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

	execMode := generator.ExecModeDirect
	if *useGoRun {
		execMode = generator.ExecModeGoRun
	} else {
		switch *flagExecMode {
		case "direct":
			execMode = generator.ExecModeDirect
		case "go_run":
			execMode = generator.ExecModeGoRun
		case "go_tool":
			execMode = generator.ExecModeGoTool
		default:
			log.Fatalf("unsupported exec mode (%s)", *flagExecMode)
		}
	}
	g := &generator.Generator{
		ExecMode:    execMode,
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
