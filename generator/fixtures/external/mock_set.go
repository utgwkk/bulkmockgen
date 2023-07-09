package external

import "database/sql/driver"

var MockInterfaces = []any{
	new(driver.Conn),
	new(driver.Driver),
}
