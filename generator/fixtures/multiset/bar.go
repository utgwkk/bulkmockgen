package multiset

import "context"

type IBar interface {
	Do(ctx context.Context) error
}
