package infrastructure

type DuplicateError struct{}

const DuplicateErrorMessage = "duplicate key error"

func (e DuplicateError) Error() string {
	return DuplicateErrorMessage
}
