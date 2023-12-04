package myerror

type MultiLoggerInterface interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Sync() error
}

type MyErrors struct {
	Text string
}

func (e *MyErrors) Error() string {
	return e.Text
}

type CheckError struct {
	RequestId int
	Logger    MultiLoggerInterface
}

type ResultError struct {
	Status  int    `json:"status"`
	Explain string `json:"explain,omitempty"`
}
