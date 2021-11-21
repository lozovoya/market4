package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"io/ioutil"
	"market4/internal/repository"
	"time"
)

const (
	PUBLICKEY = "./keys/public.key"
)

type AuthService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	usersRepo  repository.Users
	lg         *zap.Logger
}

func NewAuthService(privateKey, publicKey string, usersRepo repository.Users, lg *zap.Logger) *AuthService {
	publicKeySource, err := ioutil.ReadFile(publicKey)
	if err != nil {
		lg.Error("Auth", zap.Error(err))
		return nil
	}
	privateKeySource, err := ioutil.ReadFile(privateKey)
	if err != nil {
		lg.Error("Auth", zap.Error(err))
		return nil
	}

	k1, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeySource)
	if err != nil {
		lg.Error("Auth", zap.Error(err))
		return nil
	}
	k2, err := jwt.ParseRSAPublicKeyFromPEM(publicKeySource)
	if err != nil {
		lg.Error("Auth", zap.Error(err))
		return nil
	}

	return &AuthService{
		privateKey: k1,
		publicKey:  k2,
		usersRepo:  usersRepo,
	}
}

type Payload struct {
	ID    int
	Roles []string
	jwt.StandardClaims
}

func (a *AuthService) GetToken(id int, roles []string) (string, error) {
	payload := Payload{
		ID:    id,
		Roles: roles,
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

func (a *AuthService) GetRoleFromToken(token string) ([]string, error) {
	payload, err := jwt.ParseWithClaims(token, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		return a.publicKey, nil
	})
	var roles = make([]string, 0)
	if err != nil {
		return roles, fmt.Errorf("GetRoleFromToken: %w", err)
	}
	if !payload.Valid {
		return roles, fmt.Errorf("GetRoleFromToken: %w", err)
	}
	claims, ok := payload.Claims.(*Payload)
	if !ok {
		return roles, fmt.Errorf("GetRoleFromToken: %w", err)
	}

	return claims.Roles, nil
}

func (a *AuthService) CheckUserRole(ctx context.Context, roleID int) (string, error) {
	role, err := a.usersRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		return "", fmt.Errorf("CheckUserRole: %w", err)
	}
	return role, nil
}
