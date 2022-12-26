package apperror

import (
	"encoding/json"
	"fmt"
)

var (
	ErrNotFound = NewAppError(nil, "not found", "", "RELS-003000")
)

type ErrorFields map[string]string
type ErrorParams map[string]string

type AppError struct {
	Err              error       `json:"-"`
	Message          string      `json:"message,omitempty"`
	DeveloperMessage string      `json:"developer_message,omitempty"`
	Code             string      `json:"code,omitempty"`
	Fields           ErrorFields `json:"fields,omitempty"`
	Params           ErrorParams `json:"params,omitempty"`
}

func (e *AppError) WithFields(fields ErrorFields) {
	e.Fields = fields
}

func (e *AppError) WithParams(params ErrorParams) {
	e.Params = params
}

func NewAppError(err error, message, developerMessage, code string) *AppError {
	return &AppError{
		Err:              err,
		Message:          message,
		DeveloperMessage: developerMessage,
		Code:             code,
	}
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Marshal() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return bytes
}

func systemError(err error) *AppError {
	return NewAppError(err, "system error", err.Error(), "RELS-000001")
}

func BadRequestError(message, developerMessage string) *AppError {
	return NewAppError(fmt.Errorf(message), message, developerMessage, "RELS-000002")
}

func UnauthorizedError(message string) *AppError {
	return NewAppError(fmt.Errorf(message), message, "", "RELS-000003")
}

func APIError(message, developerMessage, code string) *AppError {
	return NewAppError(fmt.Errorf(message), message, developerMessage, code)
}
