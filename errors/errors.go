package errors

import "fmt"

type HTTPError struct {
	Status  int
	Message string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("%d - %s", e.Status, e.Message)
}
