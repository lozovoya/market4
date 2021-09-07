package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"market4/internal/api/auth"
	"market4/internal/model"
	"market4/internal/repository"
	"net/http"
)

type Auth struct {
	authService auth.AuthService
	usersRepo   repository.Users
}

func NewAuth(authService auth.AuthService, usersRepo repository.Users) *Auth {
	return &Auth{authService: authService, usersRepo: usersRepo}
}

func (a *Auth) Token(writer http.ResponseWriter, request *http.Request) {
	var data *model.User
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(fmt.Errorf("Token: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if IsEmpty(data.Login) || IsEmpty(data.Password) {
		log.Println("field is empty")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if ok := a.usersRepo.CheckCreds(request.Context(), data); !ok {
		http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	id, err := a.usersRepo.GetUserID(request.Context(), data.Login)
	if err != nil {
		log.Println(fmt.Errorf("Token: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	role, err := a.usersRepo.GetUserRole(request.Context(), data.Login)
	if err != nil {
		log.Println(fmt.Errorf("Token: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if id == 0 || role == "" {
		http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	token, err := a.authService.GetToken(request.Context(), id, role)
	if err != nil {
		log.Println(fmt.Errorf("Token: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(token)
	if err != nil {
		log.Println(fmt.Errorf("Token: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

}
