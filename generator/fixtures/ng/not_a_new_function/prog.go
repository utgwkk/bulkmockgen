package foo

var Iset = []any{
	f(2),
}

func f(x int) int {
	return x + 1
}
