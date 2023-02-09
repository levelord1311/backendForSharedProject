package apperror

type AppError string

const (
	ErrNotFound              AppError = "user not found"
	ErrCantConvertID         AppError = "can't convert ID to int"
	ErrUnpredictedInternal   AppError = "something went wrong with the server"
	ErrInvalidJSONScheme     AppError = "invalid JSON scheme. check swagger API"
	ErrAllFieldsMustBeFilled AppError = "all fields must be filled"
)

func (e AppError) Error() string {
	return string(e)
}
