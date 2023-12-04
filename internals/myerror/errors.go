package myerror

type MyErrors struct {
	Text string
}

func (e *MyErrors) Error() string {
	return e.Text
}

type ResultError struct {
	Status  int    `json:"status"`
	Explain string `json:"explain,omitempty"`
}
