package web

type ErrorResponse struct {
	Message string
	Code    uint
}

/**
Error Code
*/
const (
	NotEnoughParameters     = 10000 // Parameter
	UserIsAlreadyRegistered = 20000 // UserAccount
	UserIsNotRegistered     = 20001
	OtherError              = 90000 //Other
)
