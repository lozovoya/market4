package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

type ShopsTestSuite struct {
	suite.Suite
}

func (suite *ShopsTestSuite) SetupTest() {
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

func (suite *MarketTestSuite) TearDownTest() {
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

func Test_shopRepo_IfShopExists(t *testing.T) {
	type fields struct {
		pool *pgxpool.Pool
	}
	type args struct {
		ctx  context.Context
		shop int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &shopRepo{
				pool: tt.fields.pool,
			}
			if got := s.IfShopExists(tt.args.ctx, tt.args.shop); got != tt.want {
				t.Errorf("IfShopExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_ShopSuite(t *testing.T) {
	suite.Run(t, new(ShopsTestSuite))
}
