package apperror

type AppError string

const (
	ErrNotFound            AppError = "user not found"
	ErrCantConvertID       AppError = "can't convert ID to int"
	ErrUnpredictedInternal AppError = "something went wrong with the server"
)

func (e AppError) Error() string {
	return string(e)
}
