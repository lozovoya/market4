package md

import (
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"io/ioutil"
	"market4/internal/api/auth"
	"market4/internal/model"
	"net/http"
)

func Auth(role model.UserRole, lg *zap.Logger) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			publicKeySource, err := ioutil.ReadFile(auth.PUBLICKEY)
			if err != nil {
				lg.Error("Auth", zap.Error(err))
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeySource)
			if err != nil {
				lg.Error("Auth", zap.Error(err))
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			token := request.Header.Get("Authorization")
			if token == "" {
				lg.Error("Auth: empty token")
				http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			payload, err := jwt.ParseWithClaims(token, &auth.Payload{}, func(token *jwt.Token) (interface{}, error) {
				return publicKey, nil
			})
			if err != nil {
				lg.Error("Auth", zap.Error(err))
				http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			if !payload.Valid {
				lg.Error("Auth", zap.Error(err))
				http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			claims, ok := payload.Claims.(*auth.Payload)
			if !ok {
				http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			for _, r := range claims.Roles {
				if r == string(role) {
					handler.ServeHTTP(writer, request)
					return
				}
			}
			http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		})
	}
}
