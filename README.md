# bulkmockgen

Generate mock code all at once

## Installation

```
go install github.com/utgwkk/bulkmockgen/cmd/bulkmockgen@latest
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
$ bulkmockgen -use_go_run Iset -- -package mock_foo -destination ./mock_foo/mock.go
```

You can use bulkmockgen with `go:generate` comment.

```go
package foo

//go:generate bulkmockgen -use_go_run Iset -- -package mock_foo -destination ./mock_foo/mock.go

var Iset = []any{
	new(IFoo),
	new(IBar),
}
```

## Restriction

- Bulkmockgen is available for gomock's package mode. Source mode is currently not available.

## Migrate from mockgen

There is a migration tool `mockgen-to-bulkmockgen`. You can rewrite `//go:generate mockgen` comments to bulkmockgen's all at once.

```
go install github.com/utgwkk/bulkmockgen/cmd/mockgen-to-bulkmockgen@latest
```

Note that this migration tool can't migrate external interface mocking.
