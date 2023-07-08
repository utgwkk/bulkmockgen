package multifile

import "context"

//go:generate mockgen -package mock_foo -destination mock_foo/ifoo.go . IFoo
//go:generate mockgen -package mock_foo -destination mock_foo/ibar.go . IBar

type IFoo interface {
	Do(ctx context.Context) error
}

type IBar interface {
	Do(ctx context.Context) error
}
