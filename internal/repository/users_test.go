package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
	"log"
	"market4/internal/model"
	"sync"
	"testing"
)

type UsersTestSuite struct {
	suite.Suite
}

const (
	testDSN = "postgres://app:pass@localhost:5432/testdb"
)

func (suite *UsersTestSuite) SetupTest() {
	fmt.Println("start setup")
	testPool, err := pgxpool.Connect(context.Background(), testDSN)
	if err != nil {
		log.Println(err)
		return
	}
	createTableUsersReq := "CREATE " +
		"TABLE users ( " +
		"id BIGSERIAL PRIMARY KEY, " +
		"login TEXT NOT NULL UNIQUE, " +
		"password TEXT NOT NULL);"
	_, err = testPool.Query(context.Background(), createTableUsersReq)
	if err != nil {
		log.Println(err)
		return
	}
	createTableRolesReq := "CREATE " +
		"TABLE roles ( " +
		"id BIGSERIAL PRIMARY KEY, " +
		"name TEXT NOT NULL UNIQUE);"
	_, err = testPool.Query(context.Background(), createTableRolesReq)
	if err != nil {
		log.Println(err)
		return
	}
	createTableUserRolesReq := "CREATE " +
		"TABLE userroles ( " +
		"user_id BIGINT NOT NULL REFERENCES users, " +
		"role_id BIGINT NOT NULL REFERENCES roles, " +
		"PRIMARY KEY (user_id, role_id));"

	_, err = testPool.Query(context.Background(), createTableUserRolesReq)
	if err != nil {
		log.Println(err)
		return
	}
	addRolesReq := "INSERT " +
		"INTO roles (name) " +
		"VALUES ('USER'), ('ADMIN');"

	_, err = testPool.Query(context.Background(), addRolesReq)
	if err != nil {
		log.Println(err)
		return
	}
}

func (suite *UsersTestSuite) TearDownTest() {
	fmt.Println("cleaning up")
	testPool, err := pgxpool.Connect(context.Background(), testDSN)
	if err != nil {
		suite.Error(err)
		return
	}
	_, err = testPool.Query(context.Background(), "DROP TABLE userroles, roles, users CASCADE;")
	if err != nil {
		suite.Error(err)
	}
}

func (suite *UsersTestSuite) Test_AddUser() {
	testPool, err := pgxpool.Connect(context.Background(), testDSN)
	if err != nil {
		log.Println(err)
		return
	}
	type args struct {
		ctx  context.Context
		user *model.User
	}

	type fields struct {
		pool *pgxpool.Pool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "add user1",
			fields: fields{
				pool: testPool,
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					ID:       0,
					Login:    "user1",
					Password: "pass",
					Role:     "ADMIN",
				},
			},
			want:    &model.User{ID: 1},
			wantErr: false,
		},
		{
			name: "again add user1",
			fields: fields{
				pool: testPool,
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					ID:       0,
					Login:    "user1",
					Password: "pass",
					Role:     "ADMIN",
				},
			},
			//want:    &model.User{ID: 1},
			wantErr: true,
		},
		{
			name: "add user with wrong role",
			fields: fields{
				pool: testPool,
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					ID:       0,
					Login:    "user3",
					Password: "pass",
					Role:     "XXX",
				},
			},
			//want:    &model.User{ID: 1},
			wantErr: true,
		},
	}
	wg := sync.WaitGroup{}
	for _, tt := range tests {
		wg.Add(1)
		suite.Run(tt.name, func() {
			u := &usersRepo{
				pool: tt.fields.pool,
			}
			got, err := u.AddUser(tt.args.ctx, tt.args.user)
			fmt.Printf("GOT: %v", got)
			if (err != nil) && (tt.wantErr == true) {
				fmt.Printf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
				wg.Done()
				return
			}
			suite.Equal(tt.want, got)
			wg.Done()
		})
	}
	wg.Wait()
}

func (suite *UsersTestSuite) Test_EditUser() {

	testPool, err := pgxpool.Connect(context.Background(), testDSN)
	if err != nil {
		log.Println(err)
		return
	}

	addTestUserReq := "INSERT " +
		"INTO users (login, password) " +
		"VALUES ('user1','$2a$10$5h6GDvR0EBtCECFgptg6iuiOu0jkc/qJ8if9jt39NY9ir602nOcXu');"
	_, err = testPool.Query(context.Background(), addTestUserReq)
	if err != nil {
		suite.Error(err)
		return
	}
	type args struct {
		ctx  context.Context
		user *model.User
	}

	type fields struct {
		pool *pgxpool.Pool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "change password user1",
			fields: fields{
				pool: testPool,
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					Login:    "user1",
					Password: "qqq",
				},
			},
			want: &model.User{
				ID: 1,
			},
			wantErr: false,
		},
		{
			name: "change password wrong user",
			fields: fields{
				pool: testPool,
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					Login:    "user2",
					Password: "qqq",
				},
			},
			want:    &model.User{ID: 0},
			wantErr: false,
		},
	}
	wg := sync.WaitGroup{}
	for _, tt := range tests {
		wg.Add(1)
		suite.Run(tt.name, func() {
			u := &usersRepo{
				pool: tt.fields.pool,
			}
			got, err := u.EditUser(tt.args.ctx, tt.args.user)
			fmt.Printf("GOT %v: ", got)
			if (err != nil) && (tt.wantErr == true) {
				fmt.Printf("EditUser() error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					suite.Error(err)
				}
				wg.Done()
				return
			}
			if suite.Equal(tt.want, got) {
				fmt.Printf("EditUser() got = %v, want %v", got, tt.want)
			}
			wg.Done()
		})
	}
	wg.Wait()
}

func (suite *UsersTestSuite) Test_GetUserRolesByID() {

	testPool, err := pgxpool.Connect(context.Background(), testDSN)
	if err != nil {
		suite.Error(err)
		return
	}

	addTestUserReq := "INSERT " +
		"INTO users (login, password) " +
		"VALUES ('user1','$2a$10$5h6GDvR0EBtCECFgptg6iuiOu0jkc/qJ8if9jt39NY9ir602nOcXu');"
	_, err = testPool.Query(context.Background(), addTestUserReq)
	if err != nil {
		suite.Error(err)
		return
	}
	addUserRoleReq := "INSERT " +
		"INTO userroles (user_id, role_id) " +
		"VALUES (1, 2);"
	_, err = testPool.Query(context.Background(), addUserRoleReq)
	if err != nil {
		suite.Error(err)
		return
	}
	type args struct {
		ctx  context.Context
		user *model.User
	}

	type fields struct {
		pool *pgxpool.Pool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "check user1 role",
			fields: fields{
				pool: testPool,
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					ID: 1,
				},
			},
			want:    []string{"ADMIN"},
			wantErr: false,
		},
	}
	wg := sync.WaitGroup{}
	for _, tt := range tests {
		wg.Add(1)
		suite.Run(tt.name, func() {
			u := &usersRepo{
				pool: tt.fields.pool,
			}
			got, err := u.GetUserRolesByID(tt.args.ctx, tt.args.user.ID)
			fmt.Printf("Got: %v", got)
			if (err != nil) && (tt.wantErr == true) {
				fmt.Printf("GetUserRolesByID() error = %v, wantErr %v", err, tt.wantErr)
				wg.Done()
				return
			}
			if suite.Equal(tt.want, got) {
				fmt.Printf("GetUserRolesByID() got = %v, want %v", got, tt.want)
			}
			wg.Done()
		})
	}
	wg.Wait()
}

func Test_MarketSuite(t *testing.T) {
	suite.Run(t, new(UsersTestSuite))
}

func (suite *UsersTestSuite) Test_AddRole() {

	testPool, err := pgxpool.Connect(context.Background(), testDSN)
	if err != nil {
		suite.Error(err)
		return
	}

	addTestUserReq := "INSERT " +
		"INTO users (login, password) " +
		"VALUES ('user1','$2a$10$5h6GDvR0EBtCECFgptg6iuiOu0jkc/qJ8if9jt39NY9ir602nOcXu');"
	_, err = testPool.Query(context.Background(), addTestUserReq)
	if err != nil {
		suite.Error(err)
		return
	}
	type args struct {
		ctx  context.Context
		user *model.User
	}

	type fields struct {
		pool *pgxpool.Pool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "add existing role",
			fields: fields{
				pool: testPool,
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					Login: "user1",
					Role:  "USER",
				},
			},
			wantErr: false,
		},
		{
			name: "add wrong role",
			fields: fields{
				pool: testPool,
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					ID:    1,
					Login: "user1",
					Role:  "EDITOR",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			u := &usersRepo{
				pool: tt.fields.pool,
			}
			err := u.AddRole(tt.args.ctx, tt.args.user.Login, tt.args.user.Role)
			if err != nil {
				if tt.wantErr == true {
					fmt.Printf("AddRole() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				suite.Fail("test failed")
			}
		})
	}
}

func (suite *UsersTestSuite) Test_RemoveRole() {

	testPool, err := pgxpool.Connect(context.Background(), testDSN)
	if err != nil {
		suite.Error(err)
		return
	}

	addTestUserReq := "INSERT " +
		"INTO users (login, password) " +
		"VALUES ('user1','$2a$10$5h6GDvR0EBtCECFgptg6iuiOu0jkc/qJ8if9jt39NY9ir602nOcXu');"
	_, err = testPool.Query(context.Background(), addTestUserReq)
	if err != nil {
		suite.Error(err)
		return
	}
	addUserRoleReq := "INSERT " +
		"INTO userroles (user_id, role_id) " +
		"VALUES (1, 2);"
	_, err = testPool.Query(context.Background(), addUserRoleReq)
	if err != nil {
		suite.Error(err)
		return
	}
	type args struct {
		ctx  context.Context
		user *model.User
	}

	type fields struct {
		pool *pgxpool.Pool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "remove wrong role",
			fields: fields{
				pool: testPool,
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					Login: "user1",
					Role:  "EDITOR",
				},
			},
			wantErr: true,
		},
		{
			name: "remove ADMIN role",
			fields: fields{
				pool: testPool,
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					Login: "user1",
					Role:  "ADMIN",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			u := &usersRepo{
				pool: tt.fields.pool,
			}
			err := u.RemoveRole(tt.args.ctx, tt.args.user.Login, tt.args.user.Role)
			fmt.Printf("REMOVE ROLE %s, reply %v", tt.args.user.Role, err)
			if err != nil {
				if tt.wantErr == true {
					fmt.Printf("RemoveRole() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				suite.Fail("test failed")
			}
		})
	}
}

func (suite *UsersTestSuite) Test_CheckCreds() {

	testPool, err := pgxpool.Connect(context.Background(), testDSN)
	if err != nil {
		suite.Error(err)
		return
	}

	addTestUserReq := "INSERT " +
		"INTO users (login, password) " +
		"VALUES ('user1','$2a$10$5h6GDvR0EBtCECFgptg6iuiOu0jkc/qJ8if9jt39NY9ir602nOcXu');"
	_, err = testPool.Query(context.Background(), addTestUserReq)
	if err != nil {
		suite.Error(err)
		return
	}
	type args struct {
		ctx  context.Context
		user *model.User
	}

	type fields struct {
		pool *pgxpool.Pool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "right password",
			fields: fields{
				pool: testPool,
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					Login:    "user1",
					Password: "user1password",
				},
			},
			want: true,
		},
		{
			name: "wrong password",
			fields: fields{
				pool: testPool,
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					Login:    "user1",
					Password: "qqq",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			u := &usersRepo{
				pool: tt.fields.pool,
			}
			got := u.CheckCreds(tt.args.ctx, tt.args.user)
			fmt.Printf("GOT %v", got)
			if got != tt.want {
				fmt.Printf("CheckCreds() = %v, want %v", got, tt.want)
				suite.Fail("test failed")
			}
		})
	}
}

func (suite *UsersTestSuite) Test_GetUserID() {

	testPool, err := pgxpool.Connect(context.Background(), testDSN)
	if err != nil {
		suite.Error(err)
		return
	}

	addTestUserReq := "INSERT " +
		"INTO users (login, password) " +
		"VALUES ('user1','$2a$10$5h6GDvR0EBtCECFgptg6iuiOu0jkc/qJ8if9jt39NY9ir602nOcXu');"
	_, err = testPool.Query(context.Background(), addTestUserReq)
	if err != nil {
		suite.Error(err)
		return
	}
	type args struct {
		ctx  context.Context
		user *model.User
	}

	type fields struct {
		pool *pgxpool.Pool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "check existing user",
			fields: fields{
				pool: testPool,
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					Login: "user1",
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "check non-existing user",
			fields: fields{
				pool: testPool,
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					Login: "user2",
				},
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			u := &usersRepo{
				pool: tt.fields.pool,
			}
			got, err := u.GetUserID(tt.args.ctx, tt.args.user.Login)
			if (err != nil) != tt.wantErr {
				fmt.Printf("GetUserID() error = %v, wantErr %v", err, tt.wantErr)
				suite.Fail("test failed")
				return
			}
			if got != tt.want {
				fmt.Printf("GetUserID() got = %v, want %v", got, tt.want)
				suite.Fail("test failed")
			}
		})
	}
}

func (suite *UsersTestSuite) Test_GetRoleByID() {

	testPool, err := pgxpool.Connect(context.Background(), testDSN)
	if err != nil {
		suite.Error(err)
		return
	}
	type args struct {
		ctx  context.Context
		user *model.User
	}

	type fields struct {
		pool *pgxpool.Pool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		roleID  int
		want    string
		wantErr bool
	}{
		{
			name: "check existing role",
			fields: fields{
				pool: testPool,
			},
			args: args{
				ctx: context.Background(),
			},
			roleID:  2,
			want:    "ADMIN",
			wantErr: false,
		},
		{
			name: "check non-existing role",
			fields: fields{
				pool: testPool,
			},
			args: args{
				ctx: context.Background(),
			},
			roleID:  20,
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			u := usersRepo{
				pool: tt.fields.pool,
			}
			got, err := u.GetRoleByID(tt.args.ctx, tt.roleID)
			if (err != nil) != tt.wantErr {
				fmt.Printf("GetRoleByID() error = %v, wantErr %v", err, tt.wantErr)
				suite.Fail("test failed")
				return
			}
			if got != tt.want {
				fmt.Printf("GetRoleByID() got = %v, want %v", got, tt.want)
				suite.Fail("test failed")
			}
		})
	}
}
