package v1

import (
	"encoding/json"
	"go.uber.org/zap"
	"market4/internal/model"
	"market4/internal/repository"
	"net/http"
)

type Users struct {
	usersRepo repository.Users
	lg        *zap.Logger
}

func NewUser(usersRepo repository.Users, lg *zap.Logger) *Users {
	return &Users{usersRepo: usersRepo, lg: lg}
}

func (u *Users) AddUser(writer http.ResponseWriter, request *http.Request) {
	var data *model.User
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		u.lg.Error("AddUser", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if IsEmpty(data.Login) || IsEmpty(data.Password) || IsEmpty(data.Role) {
		u.lg.Error("AddUser: field is empty")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	addedUser, err := u.usersRepo.AddUser(request.Context(), data)
	if err != nil {
		u.lg.Error("AddUser", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(addedUser)
	if err != nil {
		u.lg.Error("AddUser", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (u *Users) EditUser(writer http.ResponseWriter, request *http.Request) {
	var data *model.User
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		u.lg.Error("EditUser", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if IsEmpty(data.Login) || IsEmpty(data.Password) {
		u.lg.Error("EditUser: field is empty")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	editedUser, err := u.usersRepo.EditUser(request.Context(), data)
	if err != nil {
		u.lg.Error("EditUser", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(editedUser)
	if err != nil {
		u.lg.Error("EditUser", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (u *Users) AddRole(writer http.ResponseWriter, request *http.Request) {
	var data *model.User
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		u.lg.Error("AddRole", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if IsEmpty(data.Login) || IsEmpty(data.Role) {
		u.lg.Error("AddRole: field is empty")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = u.usersRepo.AddRole(request.Context(), data.Login, data.Role)
	if err != nil {
		u.lg.Error("AddRole", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (u *Users) RemoveRole(writer http.ResponseWriter, request *http.Request) {
	var data *model.User
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		u.lg.Error("RemoveRole", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if IsEmpty(data.Login) || IsEmpty(data.Role) {
		u.lg.Error("RemoveRole: field is empty")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = u.usersRepo.RemoveRole(request.Context(), data.Login, data.Role)
	if err != nil {
		u.lg.Error("RemoveRole", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}
