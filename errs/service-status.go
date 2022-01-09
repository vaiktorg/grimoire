package errs

type StatusCode uint32

const (
	SUCCESS StatusCode = 0
	FAILED  StatusCode = 1
	TIMEOUT StatusCode = 2
)

type Status struct {
	Code    StatusCode
	Message string
}

func NewFailed(msg string) Status {
	return Status{Message: msg, Code: FAILED}
}
