package usecase

import (
	"atwell/domain"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type YahooJapanGetUserEmailUsecase struct {
	infra YahooJapanAuthInfrastructure
}

type YahooJapanAuthInfrastructure interface {
	GetPublicKeyList() (*http.Response, error)
	Token(code string) (*http.Response, error)
	UserInfo(accessToken string) (*http.Response, error)
}

type YahooJapanAuthenticationInformation struct {
	Token string
}

// DummyMethod is just dummy method for confirming to implement AuthenticationInformation
func (i *YahooJapanAuthenticationInformation) DummyMethod() {}

func NewYahooJapanGetUserEmailUsecase(infra YahooJapanAuthInfrastructure) domain.GetUserEmailUsecase {
	return &YahooJapanGetUserEmailUsecase{infra: infra}
}

func (u *YahooJapanGetUserEmailUsecase) GetEmail(authInfo domain.AuthenticationInformation) (email string, err error) {
	i, ok := authInfo.(*YahooJapanAuthenticationInformation)
	if !ok {
		return "", errors.New("auth information type is not YahooJapanAuthenticationInformation")
	}

	tokenData, err := u.GetToken(i.Token)
	if err != nil {
		return "", err
	}

	res, err := u.infra.UserInfo(tokenData.AccessToken)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		var errorResponse YahooJapanAPIErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return "", err
		}

		return "", fmt.Errorf(
			"failed to get user info by yahoo japan api. Error: %s, ErrorDescription: %s, ErrorCode: %d",
			errorResponse.Error,
			errorResponse.ErrorDescription,
			errorResponse.ErrorCode,
		)
	}

	var response UserInfoResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.Email, nil
}

func (u YahooJapanGetUserEmailUsecase) GetToken(code string) (TokenAPIResponse, error) {
	res, err := u.infra.Token(code)
	if err != nil {
		return TokenAPIResponse{}, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return TokenAPIResponse{}, err
	}

	if res.StatusCode != http.StatusOK {
		var tokenErrorResponse YahooJapanAPIErrorResponse
		err = json.Unmarshal(body, &tokenErrorResponse)
		if err != nil {
			return TokenAPIResponse{}, err
		}

		return TokenAPIResponse{}, fmt.Errorf(
			"failed to get ID Token. Error: %s, ErrorDescription: %s, ErrorCode: %d",
			tokenErrorResponse.Error,
			tokenErrorResponse.ErrorDescription,
			tokenErrorResponse.ErrorCode,
		)
	}

	var tokenResponse TokenAPIResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return TokenAPIResponse{}, err
	}

	// TODO: verify Token

	return tokenResponse, nil
}

func (u *YahooJapanGetUserEmailUsecase) GetPublicKey(publicKeyID string) (*rsa.PublicKey, error) {
	res, err := u.infra.GetPublicKeyList()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		var errorResponse YahooJapanAPIErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf(
			"failed to get public key. Error: %s, ErrorDescription: %s, ErrorCode: %d",
			errorResponse.Error,
			errorResponse.ErrorDescription,
			errorResponse.ErrorCode,
		)
	}

	var keys map[string]string
	err = json.Unmarshal(body, &keys)
	if err != nil {
		return nil, err
	}

	var publicKeys []PublicKey
	for keyID, keyString := range keys {
		// delete \n from keyString but remain first and last \n
		keyString = strings.Replace(keyString, "\n", "dummy for replace", 1)
		keyString = strings.Replace(keyString, "\n", "dummy for replace", -1)
		keyString = strings.ReplaceAll(keyString, "\n", "")
		keyString = strings.ReplaceAll(keyString, "dummy for replace", "\n")

		block, _ := pem.Decode([]byte(keyString))
		if block == nil {
			return nil, errors.New("failed to parse PEM block containing the public key")
		}

		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		publicKeys = append(publicKeys, PublicKey{KeyID: keyID, Key: pub.(*rsa.PublicKey)})
	}

	// extracts public pub which id match argument publicKeyID
	for _, key := range publicKeys {
		if key.KeyID == publicKeyID {
			return key.Key, err
		}
	}

	return nil, errors.New("publicKeyID does not match")
}

type PublicKey struct {
	KeyID string
	Key   *rsa.PublicKey
}

// TokenAPIResponse is struct of yahoo japan Token api response.
type TokenAPIResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
}

// UserInfoResponse is struct of yahoo japan user info api response.
type UserInfoResponse struct {
	UserID        string `json:"user_id"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

// YahooJapanAPIErrorResponse is struct of yahoo japan api response when error occurred.
type YahooJapanAPIErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorCode        int    `json:"error_code"`
}
