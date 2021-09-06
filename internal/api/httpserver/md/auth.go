package md

import (
	//"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"market4/internal/api/auth"
	"net/http"
)

func Auth(role string) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			publicKeySource, err := ioutil.ReadFile("./keys/public.key")
			if err != nil {
				log.Println(fmt.Errorf("Auth: %w", err))
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeySource)
			if err != nil {
				log.Println(fmt.Errorf("Auth: %w", err))
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			token := request.Header.Get("Authorization")
			if token == "" {
				log.Printf("Auth: empty token")
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			payload, err := jwt.ParseWithClaims(token, &auth.Payload{}, func(token *jwt.Token) (interface{}, error) {
				return publicKey, nil
			})
			if err != nil {
				log.Println(fmt.Errorf("Auth: %w", err))
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			if !payload.Valid {
				log.Println(fmt.Errorf("Auth: %w", err))
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			claims, ok := payload.Claims.(*auth.Payload)
			if !ok {
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			if claims.Role != role {
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}

			handler.ServeHTTP(writer, request)
		})
	}
}