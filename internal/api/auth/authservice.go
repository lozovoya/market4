package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"market4/internal/repository"
	"time"
)

const (
	PRIVATEKEY = "./keys/private.key"
	PUBLICKEY  = "./keys/public.key"
)

type AuthService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	usersRepo  repository.Users
}

func NewAuthService(privateKey, publicKey string, usersRepo repository.Users) *AuthService {

	publicKeySource, err := ioutil.ReadFile(publicKey)
	if err != nil {
		log.Println(fmt.Errorf("Auth: %w", err))
		return nil
	}
	privateKeySource, err := ioutil.ReadFile(privateKey)
	if err != nil {
		log.Println(fmt.Errorf("Auth: %w", err))
		return nil
	}

	k1, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeySource)
	if err != nil {
		log.Println(fmt.Errorf("Auth: %w", err))
		return nil
	}
	k2, err := jwt.ParseRSAPublicKeyFromPEM(publicKeySource)
	if err != nil {
		log.Println(fmt.Errorf("Auth: %w", err))
		return nil
	}

	return &AuthService{
		privateKey: k1,
		publicKey:  k2,
		usersRepo:  usersRepo,
	}
}

type Payload struct {
	ID   int
	Role int
	jwt.StandardClaims
}

func (a *AuthService) GetToken(ctx context.Context, id int, role int) (string, error) {

	payload := Payload{
		ID:   id,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)
	token, err := t.SignedString(a.privateKey)
	if err != nil {
		return "", fmt.Errorf("Token: %w", err)
	}
	return token, nil
}

func (a *AuthService) GetRoleFromToken(ctx context.Context, token string) (int, error) {
	payload, err := jwt.ParseWithClaims(token, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		return a.publicKey, nil
	})
	if err != nil {
		log.Println(fmt.Errorf("CheckToken: %w", err))
		return 0, err
	}
	if !payload.Valid {
		log.Println(fmt.Errorf("CheckToken: %w", err))
		return 0, err
	}
	claims, ok := payload.Claims.(*Payload)
	if !ok {
		log.Println(fmt.Errorf("CheckToken: %w", err))
		return 0, err
	}

	return claims.Role, nil
}

func (a *AuthService) CheckUserRole(ctx context.Context, roleID int) (string, error) {
	role, err := a.usersRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		fmt.Errorf("CheckUserRole: %w", err)
		return "", err
	}
	return role, nil
}
