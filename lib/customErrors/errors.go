package customErrors

import (
	"fmt"
)

type RespError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type response struct {
	Error RespError `json:"error"`
}

func (e *RespError) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Message)
}

func FormatError(err error) response {
	data, _ := err.(*RespError)

	return response{
		Error: RespError{
			Code:    data.Code,
			Message: data.Message,
			Data:    data.Data,
		},
	}
}
