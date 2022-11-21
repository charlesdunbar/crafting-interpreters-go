package main


type ReturnError struct {
	value any
}

func NewReturnError(val any) *ReturnError {
	return &ReturnError{
		value: val,
	}
}

func (e ReturnError) Error() string {
	return "Return error encountered"
}