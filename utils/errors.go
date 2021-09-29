package utils

import "fmt"

type Error string

func (e Error) Error() string {
	return string(e)
}

func HTTPError(code int) Error {
	return Error(fmt.Sprintf("http error: %v", code))
}

type Errors []error

func (e Errors) Len() int {
	return len(e)
}

func (e Errors) Error() string {
	v := ""
	for _, _e := range e {
		if _e == nil {
			continue
		}
		v += _e.Error() + "\n"
	}
	return v
}

func (e Errors) Add(err error) {
	if err == nil {
		return
	}
	e = append(e, err)
}
