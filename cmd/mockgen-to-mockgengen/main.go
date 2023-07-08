package main

import (
	"flag"
	"log"

	"github.com/utgwkk/mockgengen/migrator"
)

var (
	// useGoRun = flag.Bool("use_go_run", false, "Whether to use go run command to execute mockgen; defaults to false")
	inputDir         = flag.String("in_dir", "", "Input directory to migrate")
	fileOut          = flag.String("out", "", "File to generate migrate command; defaults to stdout")
	noOverwriteInput = flag.Bool("no_overwrite_input", false, "Do not overwrite input files")
)

func main() {
	flag.Parse()

	m := &migrator.Migrator{
		InputDir:         *inputDir,
		OutputPath:       *fileOut,
		NoOverwriteInput: *noOverwriteInput,
	}
	if err := m.Migrate(); err != nil {
		log.Fatal(err)
	}
}
