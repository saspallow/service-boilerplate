package authentication

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type AuthenticationModel struct {
	Token    string `json:"token"`
	FullName string `json:"fullName"`
}

func (o *AuthenticationModel) Bind(r *http.Request) error {
	return nil
}

type CustomerModel struct {
	ID       int64  `json:"ID"`
	FullName string `json:"fullName"`
}

func (o *CustomerModel) Bind(r *http.Request) error {
	return nil
}

type SignInResponseModel struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"` // unit: seconds
	UUID        string `json:"uuid"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
