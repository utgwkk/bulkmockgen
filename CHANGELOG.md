# CHANGELOG

## Next

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
