package infrastructure

type DuplicateError struct{}

const DuplicateErrorMessage = "duplicate key error"

func (e DuplicateError) Error() string {
	return DuplicateErrorMessage
}

type NotFoundError struct{}

const NotFoundErrorMessage = "record not found"

func (e NotFoundError) Error() string {
	return NotFoundErrorMessage
}
