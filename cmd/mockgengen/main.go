package main

import (
	"flag"
	"log"

	"github.com/utgwkk/mockgengen/generator"
)

var (
	useGoRun   = flag.Bool("use_go_run", false, "Whether to use go run command to execute mockgen; defaults to false")
	packageOut = flag.String("package", "", "Package of the generated code (by mockgen).")
	fileOut    = flag.String("out", "", "File to generate comment; defaults to stdout")
)

func main() {
	flag.Parse()

	if flag.NArg() != 2 {
		log.Fatal("Expected exactly two arguments")
	}
	sourceDir := flag.Arg(0)
	mockSetName := flag.Arg(1)

	g := &generator.Generator{
		UseGoRun:      *useGoRun,
		SourcePackage: sourceDir,
		PackageName:   *packageOut,
		MockSetName:   mockSetName,
		FileOut:       *fileOut,
	}
	if err := g.Generate(); err != nil {
		log.Fatal(err)
	}
}
