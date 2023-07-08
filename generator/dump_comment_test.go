package generator

import "testing"

func TestDumpComment(t *testing.T) {
	testcases := []struct {
		name     string
		input    *dumpCommentOption
		expected string
	}{
		{
			name: "single interface",
			input: &dumpCommentOption{
				mockgenCmd:     "go run go.uber.org/mock/mockgen",
				destination:    "mock_foo/mock_foo.go",
				packageName:    "mock_foo",
				sourcePackage:  ".",
				interfaceNames: []string{"IFoo"},
			},
			expected: "//go:generate go run go.uber.org/mock/mockgen -destination mock_foo/mock_foo.go -package mock_foo . IFoo",
		},
		{
			name: "use mockgen executable directly",
			input: &dumpCommentOption{
				mockgenCmd:     "mockgen",
				destination:    "mock_foo/mock_foo.go",
				packageName:    "mock_foo",
				sourcePackage:  ".",
				interfaceNames: []string{"IFoo"},
			},
			expected: "//go:generate mockgen -destination mock_foo/mock_foo.go -package mock_foo . IFoo",
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := dumpComment(tc.input)
			if got != tc.expected {
				t.Errorf("\nexpected: %s\ngot:      %s", tc.expected, got)
			}
		})
	}
}
