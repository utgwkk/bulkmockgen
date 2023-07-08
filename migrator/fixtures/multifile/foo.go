package multifile

import "context"

//go:generate mockgen -package mock_foo -destination mock_foo/ifoo.go . IFoo

type IFoo interface {
	Do(ctx context.Context) error
}
