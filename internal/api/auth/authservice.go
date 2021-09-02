package auth

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"time"
)

const (
	PRIVATEKEY = "./keys/private.key"
	PUBLICKEY  = "./keys/public.key"
)

type authService struct {
	privateKey string
	publicKey  string
}

func NewAuthService(pool *pgxpool.Pool) {

}

type Payload struct {
	ID   int
	Role string
	jwt.StandardClaims
}

func GetToken(ctx context.Context, id int, role string) (string, error) {
	privateKeySource, err := ioutil.ReadFile(PRIVATEKEY)
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
