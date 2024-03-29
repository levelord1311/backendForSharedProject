package apperror

type AppError string

const (
	ErrNotFound              AppError = "user not found"
	ErrCantConvertID         AppError = "can't convert ID to int"
	ErrInvalidJSONScheme     AppError = "invalid JSON scheme. check swagger API"
	ErrAllFieldsMustBeFilled AppError = "all fields must be filled"
	ErrWrongCredentials      AppError = "wrong login and/or password"
)

func (e AppError) Error() string {
	return string(e)
}
