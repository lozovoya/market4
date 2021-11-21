package v1

import (
	"encoding/json"
	"go.uber.org/zap"
	"market4/internal/api/auth"
	"market4/internal/model"
	"market4/internal/repository"
	"net/http"
)

type Token struct {
	Token string `json:"token"`
}

type Auth struct {
	authService auth.AuthService
	usersRepo   repository.Users
	lg          *zap.Logger
}

func NewAuth(authService auth.AuthService, usersRepo repository.Users, lg *zap.Logger) *Auth {
	return &Auth{authService: authService, usersRepo: usersRepo, lg: lg}
}

func (a *Auth) Token(writer http.ResponseWriter, request *http.Request) {
	var data *model.User
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		a.lg.Error("Token", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if IsEmpty(data.Login) || IsEmpty(data.Password) {
		a.lg.Error("Token: field is empty")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if ok := a.usersRepo.CheckCreds(request.Context(), *data); !ok {
		http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	id, err := a.usersRepo.GetUserID(request.Context(), data.Login)
	if err != nil {
		a.lg.Error("Token", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	roles, err := a.usersRepo.GetUserRolesByID(request.Context(), id)
	if err != nil {
		a.lg.Error("Token", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if id == 0 || len(roles) == 0 {
		http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	var reply Token
	reply.Token, err = a.authService.GetToken(id, roles)
	if err != nil {
		a.lg.Error("Token", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(reply)
	if err != nil {
		a.lg.Error("Token", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
