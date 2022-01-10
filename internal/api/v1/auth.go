package v1

import (
	"encoding/json"
	"market4/internal/api/auth"
	"market4/internal/model"
	"market4/internal/repository"
	"net/http"

	"github.com/unrolled/render"
	"go.uber.org/zap"
)

type Token struct {
	Token string `json:"token"`
}

type Auth struct {
	authService auth.AuthService
	usersRepo   repository.Users
	lg          *zap.Logger
	renderer    *render.Render
}

func NewAuth(authService auth.AuthService,
	usersRepo repository.Users,
	lg *zap.Logger,
	renderer *render.Render) *Auth {
	return &Auth{authService: authService,
		usersRepo: usersRepo,
		lg:        lg,
		renderer:  renderer}
}

func (a *Auth) Token(writer http.ResponseWriter, request *http.Request) {
	var data *model.User
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		a.lg.Error("Token", zap.Error(err))
		err = a.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "Bad request"})
		if err != nil {
			a.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	err = checkMandatoryFields(data.Login, data.Password)
	if err != nil {
		a.lg.Error("Token or password field is empty")
		err = a.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "Bad request"})
		if err != nil {
			a.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	if ok := a.usersRepo.CheckCreds(request.Context(), *data); !ok {
		err = a.renderer.JSON(writer, http.StatusUnauthorized, map[string]string{"Error": "Unauthorized"})
		if err != nil {
			a.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	id, err := a.usersRepo.GetUserID(request.Context(), data.Login)
	if err != nil {
		a.lg.Error("Token", zap.Error(err))
		err = a.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			a.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	roles, err := a.usersRepo.GetUserRolesByID(request.Context(), id)
	if err != nil {
		a.lg.Error("Token", zap.Error(err))
		err = a.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			a.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	if id == 0 || len(roles) == 0 {
		err = a.renderer.JSON(writer, http.StatusUnauthorized, map[string]string{"Error": "Unauthorized"})
		if err != nil {
			a.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	var reply Token
	reply.Token, err = a.authService.GetToken(id, roles)
	if err != nil {
		a.lg.Error("Token", zap.Error(err))
		err = a.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "Bad request"})
		if err != nil {
			a.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(reply)
	if err != nil {
		a.lg.Error("Token", zap.Error(err))
		err = a.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			a.lg.Error("Auth", zap.Error(err))
		}
		return
	}
}
