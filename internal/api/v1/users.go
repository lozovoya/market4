package v1

import (
	"encoding/json"
	"market4/internal/model"
	"market4/internal/repository"
	"net/http"

	"github.com/unrolled/render"
	"go.uber.org/zap"
)

type Users struct {
	usersRepo repository.Users
	lg        *zap.Logger
	renderer  *render.Render
}

func NewUser(usersRepo repository.Users, lg *zap.Logger, renderer *render.Render) *Users {
	return &Users{usersRepo: usersRepo, lg: lg, renderer: renderer}
}

func (u *Users) AddUser(writer http.ResponseWriter, request *http.Request) {
	var data *model.User
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		u.lg.Error("AddUser", zap.Error(err))
		err = u.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			u.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	if IsEmpty(data.Login) || IsEmpty(data.Password) || IsEmpty(data.Role) {
		u.lg.Error("AddUser: field is empty")
		err = u.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			u.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	addedUser, err := u.usersRepo.AddUser(request.Context(), data)
	if err != nil {
		u.lg.Error("AddUser", zap.Error(err))
		err = u.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			u.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(addedUser)
	if err != nil {
		u.lg.Error("AddUser", zap.Error(err))
		err = u.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			u.lg.Error("Auth", zap.Error(err))
		}
		return
	}
}

func (u *Users) EditUser(writer http.ResponseWriter, request *http.Request) {
	var data *model.User
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		u.lg.Error("EditUser", zap.Error(err))
		err = u.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			u.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	if IsEmpty(data.Login) || IsEmpty(data.Password) {
		u.lg.Error("EditUser: field is empty")
		err = u.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			u.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	editedUser, err := u.usersRepo.EditUser(request.Context(), data)
	if err != nil {
		u.lg.Error("EditUser", zap.Error(err))
		err = u.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			u.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(editedUser)
	if err != nil {
		u.lg.Error("EditUser", zap.Error(err))
		err = u.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			u.lg.Error("Auth", zap.Error(err))
		}
		return
	}
}

func (u *Users) AddRole(writer http.ResponseWriter, request *http.Request) {
	var data *model.User
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		u.lg.Error("AddRole", zap.Error(err))
		err = u.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			u.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	err = checkMandatoryFields(data.Login, data.Role)
	// if IsEmpty(data.Login) || IsEmpty(data.Role) {
	if err != nil {
		u.lg.Error("AddRole: field is empty")
		err = u.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			u.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	err = u.usersRepo.AddRole(request.Context(), data.Login, data.Role)
	if err != nil {
		u.lg.Error("AddRole", zap.Error(err))
		err = u.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			u.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (u *Users) RemoveRole(writer http.ResponseWriter, request *http.Request) {
	var data *model.User
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		u.lg.Error("RemoveRole", zap.Error(err))
		err = u.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			u.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	if IsEmpty(data.Login) || IsEmpty(data.Role) {
		u.lg.Error("RemoveRole: field is empty")
		err = u.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			u.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	err = u.usersRepo.RemoveRole(request.Context(), data.Login, data.Role)
	if err != nil {
		u.lg.Error("RemoveRole", zap.Error(err))
		err = u.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			u.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	writer.WriteHeader(http.StatusOK)
}
