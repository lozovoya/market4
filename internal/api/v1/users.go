package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"market4/internal/model"
	"market4/internal/repository"
	"net/http"
)

type Users struct {
	usersRepo repository.Users
}

func NewUser(usersRepo repository.Users) *Users {
	return &Users{usersRepo: usersRepo}
}

func (u *Users) AddUser(writer http.ResponseWriter, request *http.Request) {
	var data *model.User
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(fmt.Errorf("AddUser: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if IsEmpty(data.Login) || IsEmpty(data.Password) || IsEmpty(data.Role) {
		log.Println("field is empty")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	addedUser, err := u.usersRepo.AddUser(request.Context(), data)
	if err != nil {
		log.Println(fmt.Errorf("AddUser: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(addedUser)
	if err != nil {
		log.Println(fmt.Errorf("AddUser: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (u *Users) EditUser(writer http.ResponseWriter, request *http.Request) {
	var data *model.User
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(fmt.Errorf("EditUser: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if IsEmpty(data.Login) || IsEmpty(data.Password) {
		log.Println("field is empty")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	editedUser, err := u.usersRepo.EditUser(request.Context(), data)
	if err != nil {
		log.Println(fmt.Errorf("EditUser: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(editedUser)
	if err != nil {
		log.Println(fmt.Errorf("EditUser: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (u *Users) AddRole(writer http.ResponseWriter, request *http.Request) {
	var data *model.User
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(fmt.Errorf("AddRole: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if IsEmpty(data.Login) || IsEmpty(data.Role) {
		log.Println("field is empty")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = u.usersRepo.AddRole(request.Context(), data.Login, data.Role)
	if err != nil {
		log.Println(fmt.Errorf("AddRole: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (u *Users) RemoveRole(writer http.ResponseWriter, request *http.Request) {
	var data *model.User
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(fmt.Errorf("RemoveRole: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if IsEmpty(data.Login) || IsEmpty(data.Role) {
		log.Println("field is empty")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = u.usersRepo.RemoveRole(request.Context(), data.Login, data.Role)
	if err != nil {
		log.Println(fmt.Errorf("RemoveRole: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}
