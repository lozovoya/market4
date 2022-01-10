package repository

import (
	"context"
	"fmt"
	"market4/internal/model"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

type PricesTestSuite struct {
	suite.Suite
	testRepo  priceRepo
	productID string
	Data      TestData
}

func Test_PricesSuite(t *testing.T) {
	suite.Run(t, new(PricesTestSuite))
}

func (s *PricesTestSuite) SetupTest() {
	fmt.Println("start setup")
	var err error
	s.testRepo.pool, err = pgxpool.Connect(context.Background(), testDSN)
	if err != nil {
		s.Error(err)
		s.Fail("setup failed")
		return
	}
	s.Data, err = loadTestDataFromYaml("prices_test.yaml")
	if err != nil {
		s.Error(err)
		s.Fail("setup failed")
		return
	}
	for i, r := range s.Data.Conf.Setup.Requests {
		_, err = s.testRepo.pool.Exec(context.Background(), r.Request)
		if err != nil {
			s.Error(err)
			return
		}
		if i == 1 {
			break
		}
	}
	err = s.testRepo.pool.QueryRow(context.Background(), s.Data.Conf.Setup.Requests[2].Request).Scan(&s.productID)
	if err != nil {
		s.Error(err)
		return
	}
	addPriceReq := fmt.Sprintf(s.Data.Conf.Setup.Requests[3].Request, s.productID)
	_, err = s.testRepo.pool.Exec(context.Background(), addPriceReq)
	if err != nil {
		s.Error(err)
		s.Fail("setup failed: addPriceReq ")
		return
	}
}

func (s *PricesTestSuite) TearDownTest() {
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

func (s *PricesTestSuite) Test_priceRepo_AddPrice() {
	type args struct {
		ctx context.Context
		p   *model.Price
	}
	tests := []struct {
		name    string
		args    args
		want    model.Price
		wantErr bool
	}{
		{
			name: "price of existing product id",
			args: args{
				ctx: context.Background(),
				p: &model.Price{
					SalePrice:     10000,
					FactoryPrice:  5000,
					DiscountPrice: 7000,
					IsActive:      true,
					ProductID:     s.productID,
				},
			},
			want: model.Price{
				ID:            0,
				SalePrice:     10000,
				FactoryPrice:  5000,
				DiscountPrice: 7000,
			},
			wantErr: false,
		},
		{
			name: "price of non-existing product id",
			args: args{
				ctx: context.Background(),
				p: &model.Price{
					SalePrice:     10000,
					FactoryPrice:  5000,
					DiscountPrice: 7000,
					IsActive:      true,
					ProductID:     "111",
				},
			},
			wantErr: true,
		},
	}
	for i := range tests {
		tt := tests[i]
		s.Run(tt.name, func() {
			got, err := s.testRepo.AddPrice(tt.args.ctx, tt.args.p)
			if (err != nil) != tt.wantErr {
				fmt.Printf("AddPrice() error = %v, wantErr %v", err, tt.wantErr)
				s.Fail("AddPrice test failed", err)
				return
			}
			if !s.Equal(tt.want, got) {
				fmt.Printf("AddPrice() got = %v, want %v", got, tt.want)
				s.Fail("AddPrice test failed")
			}
		})
	}
}

func (s *PricesTestSuite) Test_priceRepo_EditPrice() {
	type args struct {
		ctx context.Context
		p   *model.Price
	}
	tests := []struct {
		name    string
		args    args
		want    model.Price
		wantErr bool
	}{
		{
			name: "edit existing price",
			args: args{
				ctx: context.Background(),
				p: &model.Price{
					ID:            1,
					SalePrice:     0,
					FactoryPrice:  0,
					DiscountPrice: 0,
					IsActive:      false,
				},
			},
			want: model.Price{
				ID:            1,
				SalePrice:     0,
				FactoryPrice:  0,
				DiscountPrice: 0,
				IsActive:      false,
				ProductID:     s.productID,
			},
			wantErr: false,
		},
		{
			name: "edit non-existing price",
			args: args{
				ctx: context.Background(),
				p: &model.Price{
					ID:            10,
					SalePrice:     0,
					FactoryPrice:  0,
					DiscountPrice: 0,
					IsActive:      false,
				},
			},
			wantErr: true,
		},
	}
	for i := range tests {
		tt := tests[i]
		s.Run(tt.name, func() {
			got, err := s.testRepo.EditPrice(tt.args.ctx, tt.args.p)
			if err != nil {
				if tt.wantErr == true {
					return
				}
				fmt.Printf("EditPrice() error = %v, wantErr %v", err, tt.wantErr)
				s.Fail("test EditPrice failed", err)
				return
			}
			if !s.Equal(tt.want, got) {
				fmt.Printf("EditPrice() got = %v, want %v", got, tt.want)
				s.Fail("test EditPrice failed")
			}
		})
	}
}

func (s *PricesTestSuite) Test_priceRepo_EditPriceByProductID() {
	type args struct {
		ctx context.Context
		p   *model.Price
	}
	tests := []struct {
		name    string
		args    args
		want    model.Price
		wantErr bool
	}{
		{
			name: "edit existing product id",
			args: args{
				ctx: context.Background(),
				p: &model.Price{

					SalePrice:     1000,
					FactoryPrice:  500,
					DiscountPrice: 900,
					IsActive:      true,
					ProductID:     s.productID,
				},
			},
			want: model.Price{
				ID:            1,
				SalePrice:     1000,
				FactoryPrice:  500,
				DiscountPrice: 900,
				IsActive:      true,
				ProductID:     s.productID,
			},
			wantErr: false,
		},
		{
			name: "edit non-existing product id",
			args: args{
				ctx: context.Background(),
				p: &model.Price{

					SalePrice:     1000,
					FactoryPrice:  500,
					DiscountPrice: 900,
					IsActive:      true,
					ProductID:     "0",
				},
			},
			wantErr: true,
		},
	}
	for i := range tests {
		tt := tests[i]
		s.Run(tt.name, func() {
			got, err := s.testRepo.EditPriceByProductID(tt.args.ctx, tt.args.p)
			if err != nil {
				if tt.wantErr == true {
					return
				}
				fmt.Printf("EditPriceByProductID() error = %v, wantErr %v", err, tt.wantErr)
				s.Fail("test EditPriceByProductID failed", err)
				return
			}
			if !s.Equal(tt.want, got) {
				fmt.Printf("EditPriceByProductID() got = %v, want %v", got, tt.want)
				s.Fail("test EditPriceByProductID failed")
			}
		})
	}
}

func (s *PricesTestSuite) Test_priceRepo_ListAllPrices() {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []model.Price
		wantErr bool
	}{
		{
			name: "request all prices",
			args: args{
				ctx: context.Background(),
			},
			want: []model.Price{
				{ID: 1,
					SalePrice:     2000,
					FactoryPrice:  1000,
					DiscountPrice: 1600,
					IsActive:      true,
					ProductID:     s.productID},
			},
			wantErr: false,
		},
	}
	for i := range tests {
		tt := tests[i]
		s.Run(tt.name, func() {
			got, err := s.testRepo.ListAllPrices(tt.args.ctx)
			if err != nil {
				if tt.wantErr == true {
					return
				}
				fmt.Printf("ListAllPrices() error = %v, wantErr %v", err, tt.wantErr)
				s.Fail("test ListAllPrices failed", err)
				return
			}
			if !s.Equal(tt.want, got) {
				fmt.Printf("ListAllPrices() got = %v, want %v", got, tt.want)
				s.Fail("test ListAllPrices failed")
			}
		})
	}
}

func (s *PricesTestSuite) Test_priceRepo_SearchPriceByProductID() {
	type args struct {
		ctx       context.Context
		productID string
	}
	tests := []struct {
		name    string
		args    args
		want    model.Price
		wantErr bool
	}{
		{
			name: "search existing price",
			args: args{
				ctx:       context.Background(),
				productID: s.productID,
			},
			want: model.Price{
				ID:            1,
				SalePrice:     2000,
				FactoryPrice:  1000,
				DiscountPrice: 1600,
				IsActive:      true,
				ProductID:     s.productID,
			},
			wantErr: false,
		},
		{
			name: "search non-existing price",
			args: args{
				ctx:       context.Background(),
				productID: "9efd8091-67ef-4c97-bb35-7cdfb1680c59",
			},
			want: model.Price{
				ID: 0,
			},
			wantErr: false,
		},
	}
	for i := range tests {
		tt := tests[i]
		s.Run(tt.name, func() {
			got, err := s.testRepo.SearchPriceByProductID(tt.args.ctx, tt.args.productID)
			if err != nil {
				if tt.wantErr == true {
					return
				}
				fmt.Printf("SearchPriceByProductID() error = %v, wantErr %v", err, tt.wantErr)
				s.Fail("test SearchPriceByProductID failed", err)
				return
			}
			if !s.Equal(got, tt.want) {
				fmt.Printf("SearchPriceByProductID() got = %v, want %v", got, tt.want)
				s.Fail("test SearchPriceByProductID failed")
			}
		})
	}
}
