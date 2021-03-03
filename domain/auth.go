package domain

type IDTokenHeader struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
	Kid string `json:"kid"`
}

type IDTokenPayload struct {
	Iss      string   `json:"iss"`
	Sub      string   `json:"sub"`
	Aud      []string `json:"aud"`
	Exp      int      `json:"exp"`
	Iat      int      `json:"iat"`
	Amr      []string `json:"amr"`
	Nonce    string   `json:"nonce"`
	AuthTime int      `json:"auth_time"`
	AtHash   string   `json:"at_hash"`
	CHash    string   `json:"c_hash"`
}

// AuthenticationInformation is a interface of information for getting user email address.
type AuthenticationInformation interface {
	// DummyMethod is just dummy method for confirming to implement AuthenticationInformation
	DummyMethod()
}

// AuthenticationUsecase is a interface for authentication
type AuthenticationUsecase interface {
	SignUp(authInfo AuthenticationInformation) error
	Login(authInfo AuthenticationInformation) (token string, err error)
}

// GetUserEmailUsecase is a interface of usecase for get user email.
type GetUserEmailUsecase interface {
	GetEmail(authInfo AuthenticationInformation) (email string, err error)
}
