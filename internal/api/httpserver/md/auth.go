package md

import (
	"io/ioutil"
	"market4/internal/api/auth"
	"market4/internal/model"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/unrolled/render"
	"go.uber.org/zap"
)

func Auth(role model.UserRole, lg *zap.Logger) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			r := render.New()
			publicKeySource, err := ioutil.ReadFile(auth.PUBLICKEY)
			if err != nil {
				lg.Error("Auth", zap.Error(err))
				err = r.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
				if err != nil {
					lg.Error("Auth", zap.Error(err))
				}
				return
			}
			publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeySource)
			if err != nil {
				lg.Error("Auth", zap.Error(err))
				err = r.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
				if err != nil {
					lg.Error("Auth", zap.Error(err))
				}
				return
			}
			token := request.Header.Get("Authorization")
			if token == "" {
				lg.Error("Auth: empty token")
				err = r.JSON(writer, http.StatusUnauthorized, map[string]string{"Error": "Unauthorized"})
				if err != nil {
					lg.Error("Auth", zap.Error(err))
				}
				return
			}
			payload, err := jwt.ParseWithClaims(token, &auth.Payload{}, func(token *jwt.Token) (interface{}, error) {
				return publicKey, nil
			})
			if err != nil {
				lg.Error("Auth", zap.Error(err))
				err = r.JSON(writer, http.StatusUnauthorized, map[string]string{"Error": "Unauthorized"})
				if err != nil {
					lg.Error("Auth", zap.Error(err))
				}
				return
			}
			if !payload.Valid {
				lg.Error("Auth", zap.Error(err))
				err = r.JSON(writer, http.StatusUnauthorized, map[string]string{"Error": "Unauthorized"})
				if err != nil {
					lg.Error("Auth", zap.Error(err))
				}
				return
			}
			claims, ok := payload.Claims.(*auth.Payload)
			if !ok {
				err = r.JSON(writer, http.StatusUnauthorized, map[string]string{"Error": "Unauthorized"})
				if err != nil {
					lg.Error("Auth", zap.Error(err))
				}
				return
			}

			for _, r := range claims.Roles {
				if r == string(role) {
					handler.ServeHTTP(writer, request)
					return
				}
			}
			err = r.JSON(writer, http.StatusUnauthorized, map[string]string{"Error": "Unauthorized"})
			if err != nil {
				lg.Error("Auth", zap.Error(err))
			}
		})
	}
}
