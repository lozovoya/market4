package repository

import (
	"context"
	"fmt"
	"market4/internal/model"
	"sync"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

type UsersTestSuite struct {
	suite.Suite
	testRepo usersRepo
	Data     TestData
}

const (
	testDSN = "postgres://app:pass@localhost:5432/testdb"
)

func (s *UsersTestSuite) SetupTest() {
	fmt.Println("start setup")
	var err error
	s.testRepo.pool, err = pgxpool.Connect(context.Background(), testDSN)
	if err != nil {
		s.Error(err)
		s.Fail("setup failed")
		return
	}
	s.Data, err = loadTestDataFromYaml("users_test.yaml")
	if err != nil {
		s.Error(err)
		s.Fail("setup failed")
		return
	}
	for _, r := range s.Data.Conf.Setup.Requests {
		_, err = s.testRepo.pool.Exec(context.Background(), r.Request)
		if err != nil {
			s.Error(err)
			return
		}
	}
}

func (s *UsersTestSuite) TearDownTest() {
	fmt.Println("cleaning up")
	var err error
	for _, r := range s.Data.Conf.Teardown.Requests {
		_, err = s.testRepo.pool.Exec(context.Background(), r.Request)
		if err != nil {
			s.Error(err)
			s.Fail("cleaning failed")
		}
	}
}

func (s *UsersTestSuite) Test_AddUser() {
	type args struct {
		ctx  context.Context
		user *model.User
	}

	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "add user1",
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
	for i := range tests {
		tt := tests[i]
		wg.Add(1)
		s.Run(tt.name, func() {
			got, err := s.testRepo.AddUser(tt.args.ctx, tt.args.user)
			fmt.Printf("GOT: %v", got)
			if (err != nil) && (tt.wantErr == true) {
				fmt.Printf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
				wg.Done()
				return
			}
			s.Equal(tt.want, got)
			wg.Done()
		})
	}
	wg.Wait()
}

func (s *UsersTestSuite) Test_EditUser() {
	addTestUserReq := "INSERT " +
		"INTO users (login, password) " +
		"VALUES ('user1','$2a$10$5h6GDvR0EBtCECFgptg6iuiOu0jkc/qJ8if9jt39NY9ir602nOcXu');"
	var err error
	_, err = s.testRepo.pool.Exec(context.Background(), addTestUserReq)
	if err != nil {
		s.Error(err)
		return
	}
	type args struct {
		ctx  context.Context
		user *model.User
	}
	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "change password user1",
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
	for i := range tests {
		tt := tests[i]
		wg.Add(1)
		s.Run(tt.name, func() {
			got, err := s.testRepo.EditUser(tt.args.ctx, tt.args.user)
			fmt.Printf("GOT %v: ", got)
			if (err != nil) && (tt.wantErr == true) {
				fmt.Printf("EditUser() error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					s.Error(err)
				}
				wg.Done()
				return
			}
			if s.Equal(tt.want, got) {
				fmt.Printf("EditUser() got = %v, want %v", got, tt.want)
			}
			wg.Done()
		})
	}
	wg.Wait()
}

func (s *UsersTestSuite) Test_GetUserRolesByID() {
	addTestUserReq := "INSERT " +
		"INTO users (login, password) " +
		"VALUES ('user1','$2a$10$5h6GDvR0EBtCECFgptg6iuiOu0jkc/qJ8if9jt39NY9ir602nOcXu');"
	var err error
	_, err = s.testRepo.pool.Exec(context.Background(), addTestUserReq)
	if err != nil {
		s.Error(err)
		return
	}
	addUserRoleReq := "INSERT " +
		"INTO userroles (user_id, role_id) " +
		"VALUES (1, 2);"
	_, err = s.testRepo.pool.Exec(context.Background(), addUserRoleReq)
	if err != nil {
		s.Error(err)
		return
	}
	type args struct {
		ctx  context.Context
		user *model.User
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "check user1 role",
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
	for i := range tests {
		tt := tests[i]
		wg.Add(1)
		s.Run(tt.name, func() {
			got, err := s.testRepo.GetUserRolesByID(tt.args.ctx, tt.args.user.ID)
			fmt.Printf("Got: %v", got)
			if (err != nil) && (tt.wantErr == true) {
				fmt.Printf("GetUserRolesByID() error = %v, wantErr %v", err, tt.wantErr)
				wg.Done()
				return
			}
			if s.Equal(tt.want, got) {
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

func (s *UsersTestSuite) Test_AddRole() {
	addTestUserReq := "INSERT " +
		"INTO users (login, password) " +
		"VALUES ('user1','$2a$10$5h6GDvR0EBtCECFgptg6iuiOu0jkc/qJ8if9jt39NY9ir602nOcXu');"
	var err error
	_, err = s.testRepo.pool.Exec(context.Background(), addTestUserReq)
	if err != nil {
		s.Error(err)
		return
	}
	type args struct {
		ctx  context.Context
		user model.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "add existing role",
			args: args{
				ctx: context.Background(),
				user: model.User{
					Login: "user1",
					Role:  "USER",
				},
			},
			wantErr: false,
		},
		{
			name: "add wrong role",
			args: args{
				ctx: context.Background(),
				user: model.User{
					ID:    1,
					Login: "user1",
					Role:  "EDITOR",
				},
			},
			wantErr: true,
		},
	}
	for i := range tests {
		tt := tests[i]
		s.Run(tt.name, func() {
			err := s.testRepo.AddRole(tt.args.ctx, tt.args.user.Login, tt.args.user.Role)
			if err != nil {
				if tt.wantErr == true {
					fmt.Printf("AddRole() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				s.Fail("test failed")
			}
		})
	}
}

func (s *UsersTestSuite) Test_RemoveRole() {
	addTestUserReq := "INSERT " +
		"INTO users (login, password) " +
		"VALUES ('user1','$2a$10$5h6GDvR0EBtCECFgptg6iuiOu0jkc/qJ8if9jt39NY9ir602nOcXu');"
	var err error
	_, err = s.testRepo.pool.Exec(context.Background(), addTestUserReq)
	if err != nil {
		s.Fail("test RemoveRole failed", err)
		return
	}
	addUserRoleReq := "INSERT " +
		"INTO userroles (user_id, role_id) " +
		"VALUES (1, 2);"
	_, err = s.testRepo.pool.Exec(context.Background(), addUserRoleReq)
	if err != nil {
		s.Fail("test RemoveRole failed", err)
		return
	}
	type args struct {
		ctx  context.Context
		user model.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "remove wrong role",
			args: args{
				ctx: context.Background(),
				user: model.User{
					Login: "user1",
					Role:  "EDITOR",
				},
			},
			wantErr: true,
		},
		{
			name: "remove ADMIN role",
			args: args{
				ctx: context.Background(),
				user: model.User{
					Login: "user1",
					Role:  "ADMIN",
				},
			},
			wantErr: false,
		},
	}
	for i := range tests {
		tt := tests[i]
		s.Run(tt.name, func() {
			err := s.testRepo.RemoveRole(tt.args.ctx, tt.args.user.Login, tt.args.user.Role)
			if err != nil {
				if tt.wantErr == true {
					return
				}
				fmt.Printf("RemoveRole() error = %v, wantErr %v", err, tt.wantErr)
				s.Fail("test RemoveRole failed")
				return
			}
		})
	}
}

func (s *UsersTestSuite) Test_CheckCreds() {
	addTestUserReq := "INSERT " +
		"INTO users (login, password) " +
		"VALUES ('user1','$2a$10$5h6GDvR0EBtCECFgptg6iuiOu0jkc/qJ8if9jt39NY9ir602nOcXu');"
	var err error
	_, err = s.testRepo.pool.Exec(context.Background(), addTestUserReq)
	if err != nil {
		s.Error(err)
		return
	}
	type args struct {
		ctx  context.Context
		user model.User
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "right password",
			args: args{
				ctx: context.Background(),
				user: model.User{
					Login:    "user1",
					Password: "user1password",
				},
			},
			want: true,
		},
		{
			name: "wrong password",
			args: args{
				ctx: context.Background(),
				user: model.User{
					Login:    "user1",
					Password: "qqq",
				},
			},
			want: false,
		},
	}
	for i := range tests {
		tt := tests[i]
		s.Run(tt.name, func() {
			got := s.testRepo.CheckCreds(tt.args.ctx, tt.args.user)
			fmt.Printf("GOT %v", got)
			if got != tt.want {
				fmt.Printf("CheckCreds() = %v, want %v", got, tt.want)
				s.Fail("test failed")
			}
		})
	}
}

func (s *UsersTestSuite) Test_GetUserID() {
	addTestUserReq := "INSERT " +
		"INTO users (login, password) " +
		"VALUES ('user1','$2a$10$5h6GDvR0EBtCECFgptg6iuiOu0jkc/qJ8if9jt39NY9ir602nOcXu');"
	var err error
	_, err = s.testRepo.pool.Exec(context.Background(), addTestUserReq)
	if err != nil {
		s.Error(err)
		return
	}
	type args struct {
		ctx  context.Context
		user *model.User
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "check existing user",
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
	for i := range tests {
		tt := tests[i]
		s.Run(tt.name, func() {
			got, err := s.testRepo.GetUserID(tt.args.ctx, tt.args.user.Login)
			if (err != nil) != tt.wantErr {
				fmt.Printf("GetUserID() error = %v, wantErr %v", err, tt.wantErr)
				s.Fail("test failed")
				return
			}
			if got != tt.want {
				fmt.Printf("GetUserID() got = %v, want %v", got, tt.want)
				s.Fail("test failed")
			}
		})
	}
}

func (s *UsersTestSuite) Test_GetRoleByID() {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		roleID  int
		want    string
		wantErr bool
	}{
		{
			name: "check existing role",
			args: args{
				ctx: context.Background(),
			},
			roleID:  2,
			want:    "ADMIN",
			wantErr: false,
		},
		{
			name: "check non-existing role",
			args: args{
				ctx: context.Background(),
			},
			roleID:  20,
			want:    "",
			wantErr: false,
		},
	}
	for i := range tests {
		tt := tests[i]
		s.Run(tt.name, func() {
			got, err := s.testRepo.GetRoleByID(tt.args.ctx, tt.roleID)
			if (err != nil) != tt.wantErr {
				fmt.Printf("GetRoleByID() error = %v, wantErr %v", err, tt.wantErr)
				s.Fail("test failed")
				return
			}
			if got != tt.want {
				fmt.Printf("GetRoleByID() got = %v, want %v", got, tt.want)
				s.Fail("test failed")
			}
		})
	}
}
