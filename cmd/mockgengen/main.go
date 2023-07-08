package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/utgwkk/mockgengen/generator"
	"golang.org/x/exp/slices"
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

	if err := g.Generate(ctx); err != nil {
		log.Fatal(err)
	}
}
