package gofaker

type CheckFailed struct {
	Message string
}

func (cf *CheckFailed) Error() string {
	return cf.Message
}
