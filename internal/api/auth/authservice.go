package auth

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"time"
)

const (
	PRIVATEKEY = "./keys/private.key"
	PUBLICKEY  = "./keys/public.key"
)

type AuthService struct {
	privateKey string
	publicKey  string
}

func NewAuthService(privateKey string) *AuthService {
	return &AuthService{
		privateKey: privateKey,
	}
}

type Payload struct {
	ID   int
	Role string
	jwt.StandardClaims
}

func (a *AuthService) GetToken(ctx context.Context, id int, role string) (string, error) {
	privateKeySource, err := ioutil.ReadFile(a.privateKey)
	if err != nil {
		return "", fmt.Errorf("Token: %w", err)
	}
	payload := Payload{
		ID:   id,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeySource)
	if err != nil {
		return "", fmt.Errorf("Token: %w", err)
	}
	token, err := t.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("Token: %w", err)
	}
	return token, nil
}
