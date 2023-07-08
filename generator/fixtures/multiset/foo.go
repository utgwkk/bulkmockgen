package multiset

import "context"

type IFoo interface {
	Do(ctx context.Context) error
}
