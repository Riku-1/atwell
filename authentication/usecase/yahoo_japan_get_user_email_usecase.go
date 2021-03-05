package usecase

import (
	"atwell/config"
	"atwell/domain"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type YahooJapanGetUserEmailUsecase struct {
	infra YahooJapanAuthInfrastructure
}

type YahooJapanAuthInfrastructure interface {
	GetPublicKeyList() (*http.Response, error)
	Token(code string) (*http.Response, error)
	UserInfo(accessToken string) (*http.Response, error)
}

func NewYahooJapanGetUserEmailUsecase(infra YahooJapanAuthInfrastructure) domain.GetUserEmailUsecase {
	return &YahooJapanGetUserEmailUsecase{infra: infra}
}

// BeforeLogin creates and returns token which contains nonce value.
func (u *YahooJapanGetUserEmailUsecase) BeforeLogin(nonce string) (token string, err error) {
	if nonce == "" {
		return "", errors.New("nonce should not be empty")
	}

	tokenJWT := jwt.New(jwt.SigningMethodHS256)
	claims := tokenJWT.Claims.(jwt.MapClaims)
	claims["yahoo_japan_nonce"] = nonce
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix() // TODO

	c, _ := config.GetAppConfig() // TODO: constructor injection
	return tokenJWT.SignedString([]byte(c.Secret))
}

func (u *YahooJapanGetUserEmailUsecase) GetEmail(code string, nonce string) (email string, err error) {

	tokenData, err := u.getToken(code)
	if err != nil {
		return "", err
	}

	err = u.verifyToken(tokenData, nonce)
	if err != nil {
		return "", err
	}

	return u.getEmailByToken(tokenData.AccessToken)
}

func (u *YahooJapanGetUserEmailUsecase) getToken(code string) (TokenAPIResponse, error) {
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
			"failed to get ID Code. Error: %s, ErrorDescription: %s, ErrorCode: %d",
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

	return tokenResponse, nil
}

func (u *YahooJapanGetUserEmailUsecase) verifyToken(tokenData TokenAPIResponse, nonce string) error {
	// verify signature
	header := strings.Split(tokenData.IDToken, ".")[0]
	decodedHeader, err := base64.RawURLEncoding.DecodeString(header)

	var idTokenHeader IDTokenHeader
	err = json.Unmarshal(decodedHeader, &idTokenHeader)

	publicKey, err := u.getPublicKey(idTokenHeader.Kid)
	if err != nil {
		return err
	}
	_, err = jwt.Parse(tokenData.IDToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			err := errors.New("unexpected signing method")
			return nil, err
		}

		return publicKey, nil
	})
	if err != nil {
		return err
	}

	// verify iss
	payload := strings.Split(tokenData.IDToken, ".")[1]
	decodedPayload, err := base64.RawURLEncoding.DecodeString(payload)
	var idTokenPayload IDTokenPayload
	err = json.Unmarshal(decodedPayload, &idTokenPayload)
	if err != nil {
		return err
	}

	if idTokenPayload.Iss != "https://auth.login.yahoo.co.jp/yconnect/v2" {
		return errors.New("iss value does not match")
	}

	// verify aud
	conf, err := config.GetYahooAuthConfig() // TODO: inject from outer
	if err != nil {
		return err
	}

	for _, id := range idTokenPayload.Aud {
		if id != conf.ClientID {
			return errors.New("aud value does not match")
		}
	}

	// verify nonce
	if idTokenPayload.Nonce != nonce {
		return errors.New("nonce does not match")
	}

	// verify access token hash
	b := sha256.Sum256([]byte(tokenData.AccessToken))
	encodedHash := base64.URLEncoding.EncodeToString(b[:len(b)/2])
	if encodedHash[:len(encodedHash)-2] != idTokenPayload.AtHash {
		return errors.New("access token hash does not match")
	}

	// verify expiration date
	if int64(idTokenPayload.Exp) < time.Now().Unix() {
		return errors.New("token is expired")
	}

	// verify iat
	const AuthenticationTime = 600 // TODO: move to another place
	if time.Now().Unix()-int64(idTokenPayload.Iat) > AuthenticationTime {
		return errors.New("token is expired")
	}

	return nil
}

func (u *YahooJapanGetUserEmailUsecase) getPublicKey(publicKeyID string) (*rsa.PublicKey, error) {
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

func (u *YahooJapanGetUserEmailUsecase) getEmailByToken(accessToken string) (email string, err error) {
	res, err := u.infra.UserInfo(accessToken)
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

type PublicKey struct {
	KeyID string
	Key   *rsa.PublicKey
}

// TokenAPIResponse is struct of yahoo japan Code api response.
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
}
