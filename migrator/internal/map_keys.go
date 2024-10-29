package internal

import (
	"cmp"
	"slices"
)

func MapKeys[M ~map[K]V, K cmp.Ordered, V any](m M) []K {
	keys := make([]K, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	slices.Sort(keys)
	return keys
}
