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
	dbReq := "INSERT INTO users (login, password) " +
		"VALUES ($1, $2) " +
		"RETURNING id"
	var addedUser model.User
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return nil, fmt.Errorf("AddUser: %w", err)
	}
	err = u.pool.QueryRow(ctx, dbReq, user.Login, hash).Scan(&addedUser.ID)
	if err != nil {
		return nil, fmt.Errorf("AddUser: %w", err)
	}

	err = u.AddRole(ctx, user.Login, user.Role)
	if err != nil {
		return nil, fmt.Errorf("AddUser: %w", err)
	}
	return &addedUser, nil
}

func (u *usersRepo) EditUser(ctx context.Context, user *model.User) (*model.User, error) {
	dbReq := "UPDATE users " +
		"SET password = $1 " +
		"WHERE login = $2 RETURNING id"
	var editedUser model.User
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return nil, fmt.Errorf("AddUser: %w", err)
	}
	err = u.pool.QueryRow(ctx, dbReq, hash, user.Login).Scan(&editedUser.ID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return &editedUser, nil
		}
		return &editedUser, fmt.Errorf("EditUser: %w", err)
	}

	return &editedUser, nil
}
func (u *usersRepo) AddRole(ctx context.Context, login, role string) error {
	dbReq := "INSERT INTO userroles (user_id, role_id) " +
		"VALUES ((SELECT id FROM users WHERE login = $1), (SELECT id FROM roles WHERE name = $2))"
	_, err := u.pool.Exec(ctx, dbReq, login, role)
	if err != nil {
		return fmt.Errorf("AddRole: %w", err)
	}
	return nil
}

func (u *usersRepo) RemoveRole(ctx context.Context, login, role string) error {
	dbReq := "DELETE FROM userroles " +
		"WHERE user_id = (SELECT id FROM users WHERE login = $1) " +
		"AND role_id = (SELECT id FROM roles WHERE name = $2)"
	_, err := u.pool.Exec(ctx, dbReq, login, role)
	if err != nil {
		return fmt.Errorf("RemoveRole: %w", err)
	}
	return nil
}

func (u *usersRepo) GetUserRolesByID(ctx context.Context, id int) ([]string, error) {
	dbReq := "SELECT role_id FROM userroles WHERE user_id = $1"
	var roles = make([]string, 0)
	rows, err := u.pool.Query(ctx, dbReq, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return roles, nil
		}
		return roles, fmt.Errorf("GetUserRoleByID: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var roleID int
		var roleName string
		err = rows.Scan(&roleID)
		if err != nil {
			return roles, fmt.Errorf("GetUserRoleByID: %w", err)
		}
		roleName, err = u.GetRoleByID(ctx, roleID)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				continue
			}
			return roles, fmt.Errorf("GetUserRoleByID: %w", err)
		}
		roles = append(roles, roleName)
	}
	return roles, nil
}

func (u *usersRepo) CheckCreds(ctx context.Context, user model.User) bool {
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
		return 0, fmt.Errorf("GetUserID: %w", err)
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
		return "", fmt.Errorf("GetRoleByID: %w", err)
	}
	return role, nil
}
