package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
	"market4/internal/model"
	"testing"
)

type ShopsTestSuite struct {
	suite.Suite
	testRepo shopRepo
}

func (suite *ShopsTestSuite) SetupTest() {
	fmt.Println("start setup")
	var err error
	suite.testRepo.pool, err = pgxpool.Connect(context.Background(), testDSN)
	if err != nil {
		suite.Error(err)
		suite.Fail("setup failed")
		return
	}
	createTableShopsReq := "CREATE TABLE shops ( " +
		"id  BIGSERIAL PRIMARY KEY, " +
		"name TEXT NOT NULL, " +
		"address TEXT NOT NULL, " +
		"lon TEXT, " +
		"lat TEXT, " +
		"working_hours   TEXT, " +
		"created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, " +
		"updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);"
	_, err = suite.testRepo.pool.Exec(context.Background(), createTableShopsReq)
	if err != nil {
		suite.Error(err)
		return
	}

	addShopReq := "INSERT " +
		"INTO shops (name, address, lon, lat, working_hours) " +
		"VALUES ('Магазин на диване', 'Москва, Останкино', '324234' , '5465476', '8 - 20');"

	_, err = suite.testRepo.pool.Exec(context.Background(), addShopReq)
	if err != nil {
		suite.Error(err)
		return
	}
}

func (suite *ShopsTestSuite) TearDownTest() {
	fmt.Println("cleaning up")
	var err error
	_, err = suite.testRepo.pool.Exec(context.Background(), "DROP TABLE shops CASCADE;")
	if err != nil {
		suite.Error(err)
	}
}

func (suite *ShopsTestSuite) Test_IfShopExists() {
	fmt.Println("start tests")
	type args struct {
		ctx  context.Context
		shop *model.Shop
	}
	fmt.Println("start tests")
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "check existing shop",
			args: args{
				ctx: context.Background(),
				shop: &model.Shop{
					ID: 1,
				},
			},
			want: true,
		},
		{
			name: "check non-existing shop",
			args: args{
				ctx: context.Background(),
				shop: &model.Shop{
					ID: 10,
				},
			},
			want: false,
		},
	}

	for i := range tests {
		tt := tests[i]
		suite.Run(tt.name, func() {
			got := suite.testRepo.IfShopExists(tt.args.ctx, tt.args.shop.ID)
			if got != tt.want {
				fmt.Printf("IfShopExists() = %v, want %v", got, tt.want)
				suite.Fail("test failed")
			}
		})
	}
}

func Test_ShopSuite(t *testing.T) {
	suite.Run(t, new(ShopsTestSuite))
}

func (suite *ShopsTestSuite) Test_ListAllShops() {
	addShopReq := "INSERT " +
		"INTO shops (name, address, lon, lat, working_hours) " +
		"VALUES ('Магазин для взрослых', 'Ростов, кремль', '12334' , '5465476', '8 - 20');"
	var err error
	_, err = suite.testRepo.pool.Exec(context.Background(), addShopReq)
	if err != nil {
		suite.Error(err)
		suite.Fail("addShopReq failed")
		return
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []model.Shop
		wantErr bool
	}{
		{
			name: "2 shops",
			args: args{
				ctx: context.Background(),
			},
			want: []model.Shop{
				{
					ID:           1,
					Name:         "Магазин на диване",
					Address:      "Москва, Останкино",
					WorkingHours: "8 - 20",
					LON:          "324234",
					LAT:          "5465476",
				},
				{
					ID:           2,
					Name:         "Магазин для взрослых",
					Address:      "Ростов, кремль",
					WorkingHours: "8 - 20",
					LON:          "12334",
					LAT:          "5465476",
				},
			},
			wantErr: false,
		},
	}
	for i := range tests {
		tt := tests[i]
		suite.Run(tt.name, func() {
			got, err := suite.testRepo.ListAllShops(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				fmt.Printf("ListAllShops() error = %v, wantErr %v", err, tt.wantErr)
				suite.Fail("test failed")
				return
			}
			if !suite.Equal(tt.want, got) {
				fmt.Printf("ListAllShops() got = %v, want %v", got, tt.want)
				suite.Fail("test failed")
			}
		})
	}
}

func (suite *ShopsTestSuite) Test_AddShop() {

	type args struct {
		ctx  context.Context
		shop *model.Shop
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "add shop",
			args: args{
				ctx: context.Background(),
				shop: &model.Shop{
					Name:         "Магазин для взрослых",
					Address:      "Ростов, кремль",
					WorkingHours: "8 - 20",
					LON:          "12334",
					LAT:          "5465476",
				},
			},
			want:    2,
			wantErr: false,
		},
	}
	for i := range tests {
		tt := tests[i]
		suite.Run(tt.name, func() {
			got, err := suite.testRepo.AddShop(tt.args.ctx, tt.args.shop)
			if (err != nil) != tt.wantErr {
				fmt.Printf("AddShop() error = %v, wantErr %v", err, tt.wantErr)
				suite.Fail("test failed")
				return
			}
			if got != tt.want {
				fmt.Printf("AddShop() got = %v, want %v", got, tt.want)
				suite.Fail("test failed")
			}
		})
	}
}

func (suite *ShopsTestSuite) Test_EditShop() {

	type args struct {
		ctx  context.Context
		shop *model.Shop
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "edit shop",
			args: args{
				ctx: context.Background(),
				shop: &model.Shop{
					ID:           1,
					Name:         "Магазин для взрослых",
					Address:      "Ростов, кремль",
					WorkingHours: "8 - 20",
					LON:          "12334",
					LAT:          "5465476",
				},
			},
			wantErr: false,
		},
	}
	for i := range tests {
		tt := tests[i]
		suite.Run(tt.name, func() {
			if err := suite.testRepo.EditShop(tt.args.ctx, tt.args.shop); (err != nil) != tt.wantErr {
				fmt.Printf("EditShop() error = %v, wantErr %v", err, tt.wantErr)
				suite.Fail("test failed")
			}
		})
	}
}
