package errors

import "fmt"

type VKErrors interface {
	GetCode() int
	GetDescription() string
}

type ApiError struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_msg"`
}

func (e *ApiError) GetCode() int {
	return e.Code
}

func (e *ApiError) GetDescription() string {
	return e.Message
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("API ERROR! Code: %d, Message: %s", e.Code, e.Message)
}
