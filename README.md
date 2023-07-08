# mockgengen
Generate go:generate comment for mockgen

## Installation

```
go install github.com/utgwkk/mockgengen/cmd/mockgengen@latest
```

## Usage

```go
// foo/foo.go
package foo

import "context"

type IFoo interface {
	Do(ctx context.Context) error
}

type IBar interface {
	Do(ctx context.Context) error
}

var Iset = []any{
	new(IFoo),
	new(IBar),
}
```

```
$ mockgengen -use_go_run Iset -- -package mock_foo -destination ./mock_foo/mock.go
```

## Restriction

- Mockgengen is available for gomock's reflect mode. Source mode is currently not available.

## Migrate from mockgen

There is a migration tool `mockgen-to-mockgengen`. You can rewrite `//go:generate mockgen` comments to mockgengen's all at once.

```
go install github.com/utgwkk/mockgengen/cmd/mockgen-to-mockgengen@latest
```
