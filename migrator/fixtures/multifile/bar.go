package multifile

import "context"

//go:generate mockgen -package mock_bar -destination mock_bar/ibar.go . IBar

type IBar interface {
	Do(ctx context.Context) error
}
