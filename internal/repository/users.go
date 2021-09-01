package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
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
	panic("implement me")
}

func (u *usersRepo) GetHash(ctx context.Context, login string) (string, error) {
	dbReq := "SELECT password FROM users WHERE login = $1"
	var hash string
	err := u.pool.QueryRow(ctx, dbReq, login).Scan(&hash)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return "", nil
		}
		fmt.Errorf("GetHash: %w", err)
		return "", err
	}
	return hash, nil
}

func (u *usersRepo) IsUserHasRole(ctx context.Context, login string, role string) (bool, error) {
	dbReq := "SELECT roles.name FROM users, roles " +
		"WHERE users.login = $1 AND users.role = roles.id"
	var reply string
	err := u.pool.QueryRow(ctx, dbReq, login).Scan(&reply)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return false, nil
		}
		fmt.Errorf("IsUserHasRole: %w", err)
		return false, err
	}
	if reply != role {
		return false, nil
	}
	return true, nil
}
