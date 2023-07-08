package generator

import (
	"strings"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestGenerate(t *testing.T) {
	testcases := []struct {
		name string
		g    *Generator
	}{
		{
			name: "multi file",
			g: &Generator{
				UseGoRun:      false,
				PackageName:   "mock_multifile",
				SourcePackage: "fixtures/multifile/",
				MockSetName:   "MockInterfaces",
			},
		},
		{
			name: "multi mock set",
			g: &Generator{
				UseGoRun:      false,
				PackageName:   "mock_multiset",
				SourcePackage: "fixtures/multiset/",
				MockSetName:   "MockFoos",
			},
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var sb strings.Builder
			tc.g.writer = &sb

			err := tc.g.Generate()
			if err != nil {
				t.Fatal(err)
			}

			snaps.MatchSnapshot(t, sb.String())
		})
	}
}
