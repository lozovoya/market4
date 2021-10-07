package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"market4/internal/model"
	"reflect"
	"testing"
)

const (
	testDSN = "postgres://app:pass@localhost:5432/testdb"
)

func Test_usersRepo_GetUserID(t *testing.T) {
	testPool, err := pgxpool.Connect(context.Background(), testDSN)
	if err != nil {
		log.Println(err)
		return
	}
	defer testPool.Close()
	createTableReq := "CREATE " +
		"TABLE users ( " +
		"id BIGSERIAL PRIMARY KEY, " +
		"login TEXT NOT NULL UNIQUE, " +
		"password TEXT NOT NULL);"
	_, err = testPool.Query(context.Background(), createTableReq)
	if err != nil {
		log.Println(err)
		return
	}
	addTestUserReq := "INSERT " +
		"INTO users (login, password) " +
		"VALUES ('user1','$2a$10$5h6GDvR0EBtCECFgptg6iuiOu0jkc/qJ8if9jt39NY9ir602nOcXu');"
	_, err = testPool.Query(context.Background(), addTestUserReq)
	if err != nil {
		log.Println(err)
		return
	}

	type fields struct {
		pool *pgxpool.Pool
	}
	type args struct {
		ctx   context.Context
		login string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "user exists",
			fields: fields{
				pool: testPool,
			},
			args: args{
				ctx:   context.Background(),
				login: "user1",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "no user",
			fields: fields{
				pool: testPool,
			},
			args: args{
				ctx:   context.Background(),
				login: "user2",
			},
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &usersRepo{
				pool: tt.fields.pool,
			}
			got, err := u.GetUserID(tt.args.ctx, tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserID() error = %v, wantErr %v", err, tt.wantErr)
				_, err = testPool.Query(context.Background(), "DROP TABLE users")
				return
			}
			if got != tt.want {
				t.Errorf("GetUserID() got = %v, want %v", got, tt.want)
			}
		})
	}
	_, err = testPool.Query(context.Background(), "DROP TABLE users")
	if err != nil {
		log.Println(err)
		return
	}

}

func Test_usersRepo_AddUser(t *testing.T) {

	testPool, err := pgxpool.Connect(context.Background(), testDSN)
	if err != nil {
		log.Println(err)
		return
	}
	defer testPool.Close()
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

	type fields struct {
		pool *pgxpool.Pool
	}
	type args struct {
		ctx  context.Context
		user *model.User
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &usersRepo{
				pool: tt.fields.pool,
			}
			got, err := u.AddUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
				_, _ = testPool.Query(context.Background(), "DROP TABLE userroles")
				_, _ = testPool.Query(context.Background(), "DROP TABLE users")
				_, _ = testPool.Query(context.Background(), "DROP TABLE roles")
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddUser() got = %v, want %v", got, tt.want)
			}
		})
	}
	_, err = testPool.Query(context.Background(), "DROP TABLE userroles")
	if err != nil {
		log.Println(err)
	}
	_, err = testPool.Query(context.Background(), "DROP TABLE users")
	if err != nil {
		log.Println(err)

	}
	_, err = testPool.Query(context.Background(), "DROP TABLE roles")
	if err != nil {
		log.Println(err)
	}
}

func Test_usersRepo_EditUser(t *testing.T) {

	testPool, err := pgxpool.Connect(context.Background(), testDSN)
	if err != nil {
		log.Println(err)
		return
	}
	defer testPool.Close()
	createTableReq := "CREATE " +
		"TABLE users ( " +
		"id BIGSERIAL PRIMARY KEY, " +
		"login TEXT NOT NULL UNIQUE, " +
		"password TEXT NOT NULL);"
	_, err = testPool.Query(context.Background(), createTableReq)
	if err != nil {
		log.Println(err)
		return
	}
	addTestUserReq := "INSERT " +
		"INTO users (login, password) " +
		"VALUES ('user1','$2a$10$5h6GDvR0EBtCECFgptg6iuiOu0jkc/qJ8if9jt39NY9ir602nOcXu');"
	_, err = testPool.Query(context.Background(), addTestUserReq)
	if err != nil {
		log.Println(err)
		return
	}

	type fields struct {
		pool *pgxpool.Pool
	}
	type args struct {
		ctx  context.Context
		user *model.User
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &usersRepo{
				pool: tt.fields.pool,
			}
			got, err := u.EditUser(tt.args.ctx, tt.args.user)
			t.Logf("GOT %v: ", got)
			if (err != nil) != tt.wantErr {
				t.Errorf("EditUser() error = %v, wantErr %v", err, tt.wantErr)
				_, err = testPool.Query(context.Background(), "DROP TABLE users")
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EditUser() got = %v, want %v", got, tt.want)
			}
		})
	}
	_, err = testPool.Query(context.Background(), "DROP TABLE users")
	if err != nil {
		log.Println(err)
		return
	}
}

func Test_usersRepo_GetUserRolesByID(t *testing.T) {

	testPool, err := pgxpool.Connect(context.Background(), testDSN)
	if err != nil {
		t.Log(err)
		return
	}
	defer testPool.Close()
	createTableUsersReq := "CREATE " +
		"TABLE users ( " +
		"id BIGSERIAL PRIMARY KEY, " +
		"login TEXT NOT NULL UNIQUE, " +
		"password TEXT NOT NULL);"
	_, err = testPool.Query(context.Background(), createTableUsersReq)
	if err != nil {
		t.Log(err)
		return
	}
	createTableRolesReq := "CREATE " +
		"TABLE roles ( " +
		"id BIGSERIAL PRIMARY KEY, " +
		"name TEXT NOT NULL UNIQUE);"
	_, err = testPool.Query(context.Background(), createTableRolesReq)
	if err != nil {
		t.Log(err)
		return
	}
	createTableUserRolesReq := "CREATE " +
		"TABLE userroles ( " +
		"user_id BIGINT NOT NULL REFERENCES users, " +
		"role_id BIGINT NOT NULL REFERENCES roles, " +
		"PRIMARY KEY (user_id, role_id));"

	_, err = testPool.Query(context.Background(), createTableUserRolesReq)
	if err != nil {
		t.Log(err)
		return
	}
	addRolesReq := "INSERT " +
		"INTO roles (name) " +
		"VALUES ('USER'), ('ADMIN');"

	_, err = testPool.Query(context.Background(), addRolesReq)
	if err != nil {
		t.Log(err)
		return
	}
	addTestUserReq := "INSERT " +
		"INTO users (login, password) " +
		"VALUES ('user1','$2a$10$5h6GDvR0EBtCECFgptg6iuiOu0jkc/qJ8if9jt39NY9ir602nOcXu');"
	_, err = testPool.Query(context.Background(), addTestUserReq)
	if err != nil {
		t.Log(err)
		return
	}

	addUserRoleReq := "INSERT " +
		"INTO userroles (user_id, role_id) " +
		"VALUES (1, 2);"
	_, err = testPool.Query(context.Background(), addUserRoleReq)
	if err != nil {
		t.Log(err)
		return
	}

	type fields struct {
		pool *pgxpool.Pool
	}
	type args struct {
		ctx context.Context
		id  int
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
				id:  1,
			},
			want:    []string{"ADMIN"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &usersRepo{
				pool: tt.fields.pool,
			}
			got, err := u.GetUserRolesByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserRolesByID() error = %v, wantErr %v", err, tt.wantErr)
				_, _ = testPool.Query(context.Background(), "DROP TABLE userroles")
				_, _ = testPool.Query(context.Background(), "DROP TABLE users")
				_, _ = testPool.Query(context.Background(), "DROP TABLE roles")
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserRolesByID() got = %v, want %v", got, tt.want)
			}
		})
	}
	_, _ = testPool.Query(context.Background(), "DROP TABLE userroles")
	_, _ = testPool.Query(context.Background(), "DROP TABLE users")
	_, _ = testPool.Query(context.Background(), "DROP TABLE roles")
}
