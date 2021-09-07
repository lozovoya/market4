package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
	"market4/internal/model"
	"strings"
)

type usersRepo struct {
	pool *pgxpool.Pool
}

func NewUsersRepo(pool *pgxpool.Pool) Users {
	return &usersRepo{pool: pool}
}

func (u *usersRepo) AddUser(ctx context.Context, user *model.User) (*model.User, error) {
	dbReq := "INSERT INTO users (login, password, role) " +
		"VALUES ($1, $2, (SELECT id FROM roles WHERE name = $3)) " +
		"RETURNING id"
	var addedUser model.User
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		fmt.Errorf("AddUser: %w", err)
		return nil, err
	}
	err = u.pool.QueryRow(ctx, dbReq, user.Login, hash, user.Role).Scan(&user.ID)
	if err != nil {
		fmt.Errorf("AddUser: %w", err)
		return nil, err
	}
	return &addedUser, nil
}

func (u *usersRepo) EditUser(ctx context.Context, user *model.User) (*model.User, error) {
	dbReq := "UPDATE users " +
		"SET password = $1, " +
		"role = (SELECT id FROM roles WHERE name = $2) " +
		"WHERE login = $3 RETURNING id"
	var editedUser model.User
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		fmt.Errorf("AddUser: %w", err)
		return nil, err
	}
	err = u.pool.QueryRow(ctx, dbReq, hash, user.Role, user.Login).Scan(&editedUser.ID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return &editedUser, nil
		}
		fmt.Errorf("EditUser: %w", err)
		return &editedUser, err
	}
	return &editedUser, nil
}

func (u *usersRepo) GetUserRole(ctx context.Context, login string) (string, error) {
	dbReq := "SELECT roles.name FROM users, roles " +
		"WHERE users.login = $1 AND users.role = roles.id"
	var role string
	err := u.pool.QueryRow(ctx, dbReq, login).Scan(&role)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return "", nil
		}
		log.Println(fmt.Errorf("GetUserRole: %w", err))
		return "", err
	}
	return role, nil
}

func (u *usersRepo) CheckCreds(ctx context.Context, user *model.User) bool {
	dbReq := "SELECT password FROM users WHERE login = $1"
	var hash []byte
	err := u.pool.QueryRow(ctx, dbReq, user.Login).Scan(&hash)
	if err != nil {
		log.Println(fmt.Errorf("CheckCreds: %w", err))
		return false
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(user.Password))
	if err != nil {
		log.Println(fmt.Errorf("CheckCreds: %w: ", err))
		return false
	}

	return true
}

func (u *usersRepo) GetUserID(ctx context.Context, login string) (int, error) {
	dbReq := "SELECT id FROM users WHERE login = $1"
	var id int
	err := u.pool.QueryRow(ctx, dbReq, login).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return 0, nil
		}
		fmt.Errorf("GetUserID: %w", err)
		return 0, err
	}
	return id, nil
}

func (u usersRepo) GetRoleByID(ctx context.Context, roleID int) (string, error) {
	dbReq := "SELECT name FROM roles WHERE id = $1"
	var role string
	err := u.pool.QueryRow(ctx, dbReq, roleID).Scan(&role)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return "", nil
		}
		fmt.Errorf("GetRoleByID: %w", err)
		return "", err
	}
	return role, nil
}
