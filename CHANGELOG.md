# CHANGELOG

## [v0.4.2](https://github.com/utgwkk/bulkmockgen/compare/v0.4.1...v0.4.2) - 2025-08-12
- build(deps): bump github.com/stoewer/go-strcase from 1.3.0 to 1.3.1 by @dependabot[bot] in https://github.com/utgwkk/bulkmockgen/pull/41
- build(deps): bump golang.org/x/tools from 0.30.0 to 0.34.0 by @dependabot[bot] in https://github.com/utgwkk/bulkmockgen/pull/39
- build(deps): bump github.com/gkampitakis/go-snaps from 0.5.11 to 0.5.13 by @dependabot[bot] in https://github.com/utgwkk/bulkmockgen/pull/40
- build(deps): bump github.com/gkampitakis/go-snaps from 0.5.13 to 0.5.14 by @dependabot[bot] in https://github.com/utgwkk/bulkmockgen/pull/43
- build(deps): bump golang.org/x/tools from 0.35.0 to 0.36.0 by @dependabot[bot] in https://github.com/utgwkk/bulkmockgen/pull/45
- build(deps): bump actions/checkout from 4 to 5 by @dependabot[bot] in https://github.com/utgwkk/bulkmockgen/pull/44

## [v0.4.1](https://github.com/utgwkk/bulkmockgen/compare/v0.4.0...v0.4.1) - 2025-04-30
- build(deps): bump go.uber.org/mock from 0.5.0 to 0.5.1 by @dependabot in https://github.com/utgwkk/bulkmockgen/pull/33
- build(deps): bump go.uber.org/mock from 0.5.1 to 0.5.2 by @dependabot in https://github.com/utgwkk/bulkmockgen/pull/36

## [v0.4.0](https://github.com/utgwkk/bulkmockgen/compare/v0.3.1...v0.4.0) - 2025-02-20
- build(deps): bump github.com/gkampitakis/go-snaps from 0.5.9 to 0.5.10 by @dependabot in https://github.com/utgwkk/bulkmockgen/pull/27
- build(deps): bump golang.org/x/tools from 0.29.0 to 0.30.0 by @dependabot in https://github.com/utgwkk/bulkmockgen/pull/28
- build(deps): bump github.com/gkampitakis/go-snaps from 0.5.10 to 0.5.11 by @dependabot in https://github.com/utgwkk/bulkmockgen/pull/30
- Support `go tool` command and deprecate `-use_go_run` argument by @utgwkk in https://github.com/utgwkk/bulkmockgen/pull/31

## [v0.3.1](https://github.com/utgwkk/bulkmockgen/compare/v0.3.0...v0.3.1) - 2025-02-04
- Remove golang.org/x/exp dependency by @utgwkk in https://github.com/utgwkk/bulkmockgen/pull/14
- deps: configure GitHub Actions dependabot update by @utgwkk in https://github.com/utgwkk/bulkmockgen/pull/15
- build(deps): bump actions/checkout from 3 to 4 by @dependabot in https://github.com/utgwkk/bulkmockgen/pull/16
- build(deps): bump actions/setup-go from 4 to 5 by @dependabot in https://github.com/utgwkk/bulkmockgen/pull/17
- doc: reflect mode is deprecated and package mode is introduced by @utgwkk in https://github.com/utgwkk/bulkmockgen/pull/18
- build(deps): bump golang.org/x/tools from 0.26.0 to 0.27.0 by @dependabot in https://github.com/utgwkk/bulkmockgen/pull/19
- build(deps): bump github.com/stretchr/testify from 1.9.0 to 1.10.0 by @dependabot in https://github.com/utgwkk/bulkmockgen/pull/20
- build(deps): bump golang.org/x/tools from 0.27.0 to 0.28.0 by @dependabot in https://github.com/utgwkk/bulkmockgen/pull/21
- build(deps): bump github.com/gkampitakis/go-snaps from 0.5.7 to 0.5.8 by @dependabot in https://github.com/utgwkk/bulkmockgen/pull/22
- build(deps): bump golang.org/x/tools from 0.28.0 to 0.29.0 by @dependabot in https://github.com/utgwkk/bulkmockgen/pull/23
- build(deps): bump github.com/gkampitakis/go-snaps from 0.5.8 to 0.5.9 by @dependabot in https://github.com/utgwkk/bulkmockgen/pull/24
- tagpr by @utgwkk in https://github.com/utgwkk/bulkmockgen/pull/25

## Version 0.3.0

- build(deps): bump github.com/stretchr/testify from 1.8.4 to 1.9.0 ([#8](https://github.com/utgwkk/bulkmockgen/pull/8))
- build(deps): bump golang.org/x/tools from 0.11.0 to 0.26.0 ([#9](https://github.com/utgwkk/bulkmockgen/pull/9))
- build(deps): bump github.com/gkampitakis/go-snaps from 0.4.8 to 0.5.7 ([#10](https://github.com/utgwkk/bulkmockgen/pull/10))
- build(deps): bump go.uber.org/mock from 0.2.0 to 0.5.0 ([#11](https://github.com/utgwkk/bulkmockgen/pull/11))
- Use go 1.20 ([#12](https://github.com/utgwkk/bulkmockgen/pull/12))
- Use go 1.22 ([#13](https://github.com/utgwkk/bulkmockgen/pull/13))
- ci: disable Go dependency cache

## Version 0.2.2

- add mockgen to tools dependencies ([#5](https://github.com/utgwkk/bulkmockgen/pull/5))

## Version 0.2.1

- Allow external package mocking like below:

```go
package external

import "database/sql/driver"

// The below command is equivalent to:
//   mockgen -package mock_driver -destination ./mock_driver/mock.go database/sql/driver Conn,Driver

//go:generate bulkmockgen MockInterfaces -- -package mock_driver -destination ./mock_driver/mock.go
var MockInterfaces = []any{
	new(driver.Conn),
	new(driver.Driver),
}
```

- You can't still mock mixed external packages' interfaces at once. Please split mock set and generatel one by one.

```go
package mixedexternal

import (
	"database/sql"
	"database/sql/driver"
)

// NG
var Iset = []any{
	new(sql.Result),
	new(driver.Conn),
}

// OK
var SqlSet = []any{
	new(sql.Result),
}

var DriverSet = []any{
	new(driver.Conn),
}

```

## Version 0.2.0 (2023/7/9)

- Rename mockgengen to **bulkmockgen**

## Version 0.1.0 (2023/7/8)

- **incompatible**: switch to wrap mockgen command
  - You can use mockgengen with mockgen's command line options.
  - eg. `mockgengen MockBars -- -package mock_bar -destination ./mock_bar/mock_bar.go`

## Version 0.0.2 (2023/7/8)

- migrator: treat consecutive go:generate comment correctly

## Version 0.0.1 (2023/7/8)

- First release
