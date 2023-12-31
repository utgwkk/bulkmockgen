package migrator

import (
	"strings"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestGenerate(t *testing.T) {
	testcases := []struct {
		name string
		m    *Migrator
	}{
		{
			name: "multi file",
			m: &Migrator{
				InputDir: "./fixtures/multifile",
			},
		},
		{
			name: "multi line",
			m: &Migrator{
				InputDir: "./fixtures/multiline",
			},
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var sb strings.Builder
			tc.m.writer = &sb
			tc.m.NoOverwriteInput = true

			err := tc.m.Migrate()
			if err != nil {
				t.Fatal(err)
			}

			snaps.MatchSnapshot(t, sb.String())
		})
	}
}
