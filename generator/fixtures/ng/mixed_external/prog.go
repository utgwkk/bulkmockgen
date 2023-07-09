package mixedexternal

import (
	"database/sql"
	"database/sql/driver"
)

var Iset = []any{
	new(sql.Result),
	new(driver.Conn),
}
