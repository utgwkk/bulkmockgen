package migrator

import (
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/stoewer/go-strcase"
	"golang.org/x/exp/maps"
)

var mockgenCommandCandidates = []string{
	"mockgen",
	"go run github.com/golang/mock/mockgen",
	"go run go.uber.org/mock/mockgen",
}

type Migrator struct {
	InputDir         string
	OutputPath       string
	NoOverwriteInput bool

	writer io.Writer // for testing
}

var (
	regexpSpaces = regexp.MustCompile(`\s+`)

	pluralizeClient = pluralize.NewClient()
)

func (m *Migrator) Migrate() error {
	fset := token.NewFileSet()
	parsed, err := parser.ParseDir(fset, m.InputDir, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	var comments []*mockgenComment
	for _, pkg := range parsed {
		for _, f := range pkg.Files {
			cs := findMockgenGoGenerateComments(f)
			if len(cs) > 0 {
				comments = append(comments, cs...)
			}
		}
	}

	if len(comments) == 0 {
		return nil
	}

	abspath, err := filepath.Abs(m.InputDir)
	if err != nil {
		return err
	}

	splittedPath := strings.Split(filepath.ToSlash(abspath), "/")
	sourcePkg := splittedPath[len(splittedPath)-1]

	gomockOpts := collectGomockCommandOpts(comments)
	type outInterfaces struct {
		sourcePkg  string
		interfaces []string
	}
	interfacesByPackageName := make(map[string]*outInterfaces)
	for _, o := range gomockOpts {
		if o.packageName == "" {
			o.packageName = "mock_" + sourcePkg
		}
		if interfacesByPackageName[o.packageName] == nil {
			interfacesByPackageName[o.packageName] = &outInterfaces{
				sourcePkg: o.sourcePkg,
			}
		}
		interfacesByPackageName[o.packageName].interfaces = append(interfacesByPackageName[o.packageName].interfaces, o.targetInterfaces...)
	}

	var w io.Writer
	if m.writer != nil {
		w = m.writer
	} else if m.OutputPath != "" {
		f, err := os.Create(m.OutputPath)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	} else {
		w = os.Stdout
	}

	// output code with determinate order
	keys := maps.Keys(interfacesByPackageName)
	sort.Strings(keys)

	w.Write([]byte(fmt.Sprintf("package %s\n\n", sourcePkg)))
	for _, pkgName := range keys {
		is := interfacesByPackageName[pkgName]
		varName := pluralizeClient.Plural(strcase.UpperCamelCase(pkgName))
		w.Write([]byte(fmt.Sprintf("//go:generate go run github.com/utgwkk/mockgengen/cmd/mockgengen %s -- -package %s -destination ./%s/%s.go\n", varName, pkgName, pkgName, pkgName)))
		w.Write([]byte(fmt.Sprintf("var %s = []any{\n", varName)))
		sort.Strings(is.interfaces)
		for _, i := range is.interfaces {
			w.Write([]byte(fmt.Sprintf("\tnew(%s),\n", i)))
		}
		w.Write([]byte("}\n"))
	}

	if m.NoOverwriteInput {
		return nil
	}

	for _, pkg := range parsed {
		for filename, f := range pkg.Files {
			f.Comments = removeGoGenerateComment(f.Comments, comments)

			err := func() error {
				fout, err := os.CreateTemp(os.TempDir(), "mockgen-to-mockgengen")
				if err != nil {
					return err
				}
				defer fout.Close()

				if err := format.Node(fout, fset, f); err != nil {
					return err
				}
				fout.Seek(0, 0)

				srcFile, err := os.Create(filename)
				if err != nil {
					return err
				}
				defer srcFile.Close()

				if _, err := io.Copy(srcFile, fout); err != nil {
					return err
				}

				return nil
			}()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func removeGoGenerateComment(comments []*ast.CommentGroup, goGenerateComments []*mockgenComment) []*ast.CommentGroup {
	var removedComments []*ast.CommentGroup
	for _, c := range comments {
		var clist []*ast.Comment
	SKIP_COMMENT:
		for _, comment := range c.List {
			for _, gc := range goGenerateComments {
				if comment == gc.comment {
					continue SKIP_COMMENT
				}
			}
			clist = append(clist, comment)
		}
		c.List = clist
		removedComments = append(removedComments, c)
	}
	return removedComments
}

type mockgenComment struct {
	comment        *ast.Comment
	normalizedText string
}

func findMockgenGoGenerateComments(f *ast.File) []*mockgenComment {
	var comments []*mockgenComment

	for _, c := range f.Comments {
		for _, comment := range c.List {
			commentText := comment.Text
			if !strings.HasPrefix(commentText, "//go:generate") {
				continue
			}

			trimmed := strings.TrimPrefix(commentText, "//go:generate ")
			normalized, ok := detectAndNormalizeMockgenGoGenerateComment(trimmed)
			if !ok {
				continue
			}

			comments = append(comments, &mockgenComment{
				comment:        comment,
				normalizedText: normalized,
			})
		}
	}
	return comments
}

func detectAndNormalizeMockgenGoGenerateComment(comment string) (string, bool) {
	for _, candid := range mockgenCommandCandidates {
		if strings.HasPrefix(comment, candid+" ") {
			normalized := strings.TrimPrefix(comment, candid+" ")
			normalized = strings.TrimSpace(normalized)
			return normalized, true
		}
	}
	return "", false
}

type gomockCommandOpts struct {
	source           string
	destination      string
	packageName      string
	sourcePkg        string
	targetInterfaces []string
}

func collectGomockCommandOpts(comments []*mockgenComment) []*gomockCommandOpts {
	var gomockOpts []*gomockCommandOpts

	for _, c := range comments {
		flagSet := flag.NewFlagSet("", flag.ContinueOnError)

		// ref: https://github.com/uber/mock/blob/dac455047760bb7061f57f42615cacfa1fac75c1/mockgen/mockgen.go#L57-L64
		flagSet.String("source", "", "(source mode) Input Go source file; enables source mode.")
		flagSet.String("destination", "", "Output file; defaults to stdout.")
		// flagSet.String("mock_names", "", "Comma-separated interfaceName=mockName pairs of explicit mock names to use. Mock names default to 'Mock'+ interfaceName suffix.")
		flagSet.String("package", "", "Package of the generated code; defaults to the package of the input with a 'mock_' prefix.")
		// flagSet.String("self_package", "", "The full package import path for the generated code. The purpose of this flag is to prevent import cycles in the generated code by trying to include its own package. This can happen if the mock's package is set to one of its inputs (usually the main one) and the output is stdio so mockgen cannot detect the final output package. Setting this flag will then tell mockgen which import to exclude.")
		// flagSet.Bool("write_package_comment", true, "Writes package documentation comment (godoc) if true.")
		// flagSet.String("copyright_file", "", "Copyright file used to add copyright header")
		// flagSet.Bool("typed", false, "Generate Type-safe 'Return', 'Do', 'DoAndReturn' function")

		args := regexpSpaces.Split(c.normalizedText, -1)
		err := flagSet.Parse(args)
		if err != nil {
			continue
		}

		flg := &gomockCommandOpts{}
		flagSet.VisitAll(func(f *flag.Flag) {
			switch f.Name {
			case "source":
				flg.source = f.Value.String()
			case "destination":
				flg.destination = f.Value.String()
			case "package":
				flg.packageName = f.Value.String()
			}
		})
		if flg.source != "" {
			// source mode is not available
			continue
		}

		sourcePkg := flagSet.Arg(0)
		flg.sourcePkg = sourcePkg

		targetInterfacesStr := flagSet.Arg(1)
		if targetInterfacesStr == "" {
			continue
		}

		targetInterfaces := strings.Split(targetInterfacesStr, ",")
		flg.targetInterfaces = targetInterfaces
		gomockOpts = append(gomockOpts, flg)
	}

	return gomockOpts
}
