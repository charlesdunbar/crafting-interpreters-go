package functiontype

type FunctionType int64

const (
	NONE FunctionType = iota
	FUNCTION
	INITIALIZER
	METHOD
)
