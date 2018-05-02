package gofaker

type AssertionFailedError struct {
	Message string
}

func (cf *AssertionFailedError) Error() string {
	return cf.Message
}
