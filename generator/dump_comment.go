package generator

import (
	"fmt"
	"strings"
)

type dumpCommentOption struct {
	mockgenCmd     string
	destination    string
	packageName    string
	sourcePackage  string
	interfaceNames []string
}

func dumpComment(opt *dumpCommentOption) string {
	joinedInterfaces := strings.Join(opt.interfaceNames, ",")
	return fmt.Sprintf(
		"//go:generate %s -destination %s -package %s %s %s",
		opt.mockgenCmd,
		opt.destination,
		opt.packageName,
		opt.sourcePackage,
		joinedInterfaces,
	)
}
