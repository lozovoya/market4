package repository

import (
	"context"
	"fmt"
	"market4/internal/model"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

type ProductTestSuite struct {
	suite.Suite
	testRepo  productRepo
	productID string
	Data      TestData
}

func Test_ProductSuite(t *testing.T) {
	suite.Run(t, new(ProductTestSuite))
}

func (s *ProductTestSuite) SetupTest() {
	fmt.Println("start setup")
	var err error

	s.testRepo.pool, err = pgxpool.Connect(context.Background(), testDSN)
	if err != nil {
		s.Fail("setup failed", err)
		return
	}
	s.Data, err = loadTestDataFromYaml("products_test.yaml")
	if err != nil {
		s.Error(err)
		s.Fail("setup failed")
		return
	}
	createExtensionReq := "CREATE EXTENSION pgcrypto;"
	_, err = s.testRepo.pool.Exec(context.Background(), createExtensionReq)
	if err != nil {
		fmt.Println("pgcrypto failed: createExtensionReq", err)
	}
	for i, r := range s.Data.Conf.Setup.Requests {
		_, err = s.testRepo.pool.Exec(context.Background(), r.Request)
		if err != nil {
			s.Error(err)
			return
		}
		if i == 6 {
			break
		}
	}
	err = s.testRepo.pool.QueryRow(context.Background(), s.Data.Conf.Setup.Requests[7].Request).Scan(&s.productID)
	if err != nil {
		s.Fail("setup failed: addProductReq", err)
		return
	}
	addProductCategoryReq := fmt.Sprintf(s.Data.Conf.Setup.Requests[8].Request, s.productID)
	_, err = s.testRepo.pool.Exec(context.Background(), addProductCategoryReq)
	if err != nil {
		s.Fail("setup failed", err)
		return
	}
	addProductShopReq := fmt.Sprintf(s.Data.Conf.Setup.Requests[9].Request, s.productID)
	_, err = s.testRepo.pool.Exec(context.Background(), addProductShopReq)
	if err != nil {
		s.Fail("setup failed", err)
		return
	}
}

func (s *ProductTestSuite) TearDownTest() {
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

func (s *ProductTestSuite) Test_productRepo_IfProductExists() {
	type args struct {
		ctx       context.Context
		productID string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "check existing product",
			args: args{
				ctx:       context.Background(),
				productID: s.productID,
			},
			want: true,
		},
		{
			name: "check non-existing product",
			args: args{
				ctx:       context.Background(),
				productID: "9efd8091-67ef-4c97-bb35-7cdfb1680c59",
			},
			want: false,
		},
	}
	for i := range tests {
		tt := tests[i]
		s.Run(tt.name, func() {
			if got := s.testRepo.IfProductExists(tt.args.ctx, tt.args.productID); got != tt.want {
				fmt.Printf("IfProductExists() = %v, want %v", got, tt.want)
				s.Fail("test IfProductExists failed")
			}
		})
	}
}

func (s *ProductTestSuite) Test_productRepo_setProductCategory() {
	type args struct {
		ctx        context.Context
		categoryId int
		productId  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "set existing category",
			args: args{
				ctx:        context.Background(),
				categoryId: 2,
				productId:  s.productID,
			},
			wantErr: false,
		},
		{
			name: "set non-existing category",
			args: args{
				ctx:        context.Background(),
				categoryId: 20,
				productId:  s.productID,
			},
			wantErr: true,
		},
	}
	for i := range tests {
		tt := tests[i]
		s.Run(tt.name, func() {
			err := s.testRepo.setProductCategory(tt.args.ctx, tt.args.categoryId, tt.args.productId)
			if err != nil {
				if tt.wantErr == true {
					return
				}
				fmt.Printf("setProductCategory() error = %v, wantErr %v", err, tt.wantErr)
				s.Fail("test setProductCategory failed")
			}
		})
	}
}

func (s *ProductTestSuite) Test_productRepo_setProductShop() {
	type args struct {
		ctx       context.Context
		shopID    int
		productID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "set existing shop",
			args: args{
				ctx:       context.Background(),
				shopID:    2,
				productID: s.productID,
			},
			wantErr: false,
		},
		{
			name: "set non-existing shop",
			args: args{
				ctx:       context.Background(),
				shopID:    10,
				productID: s.productID,
			},
			wantErr: true,
		},
	}
	for i := range tests {
		tt := tests[i]
		s.Run(tt.name, func() {
			err := s.testRepo.setProductShop(tt.args.ctx, tt.args.shopID, tt.args.productID)
			if err != nil {
				if tt.wantErr == true {
					return
				}
				fmt.Printf("setProductShop() error = %v, wantErr %v", err, tt.wantErr)
				s.Fail("test setProductShop failed")
			}
		})
	}
}

func (s *ProductTestSuite) Test_productRepo_EditProduct() {
	type args struct {
		ctx        context.Context
		product    model.Product
		shopID     int
		categoryID int
	}
	tests := []struct {
		name    string
		args    args
		want    model.Product
		wantErr bool
	}{
		{
			name: "edit product existing sku",
			args: args{
				ctx: context.Background(),
				product: model.Product{
					SKU:      "3001",
					Name:     "клюшка",
					IsActive: true,
				},
			},
			want: model.Product{
				ID:   s.productID,
				SKU:  "3001",
				Name: "клюшка",
				//Type:        "",
				Description: "пушка детская",
				URI:         "/product/тепловая-3001",
				IsActive:    true,
			},
			wantErr: false,
		},
		{
			name: "edit product non-existing sku",
			args: args{
				ctx: context.Background(),
				product: model.Product{
					SKU:      "0001",
					Name:     "клюшка",
					IsActive: true,
				},
			},
			wantErr: true,
		},
	}
	for i := range tests {
		tt := tests[i]
		s.Run(tt.name, func() {
			got, err := s.testRepo.EditProduct(tt.args.ctx, tt.args.product, tt.args.shopID, tt.args.categoryID)
			if err != nil {
				if tt.wantErr == true {
					return
				}
				fmt.Printf("EditProduct() error = %v, wantErr %v", err, tt.wantErr)
				s.Fail("test EditProduct failed")
				return
			}
			if !s.Equal(tt.want, got) {
				fmt.Printf("EditProduct() got = %v, want %v", got, tt.want)
				s.Fail("test EditProduct failed")
			}
		})
	}
}

func (s *ProductTestSuite) Test_productRepo_ListAllProducts() {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []model.Product
		wantErr bool
	}{
		{
			name: "list all products",
			args: args{
				ctx: context.Background(),
			},
			want: []model.Product{
				{
					ID:          s.productID,
					SKU:         "3001",
					Name:        "пушка",
					Type:        "",
					URI:         "/product/тепловая-3001",
					Description: "пушка детская",
					IsActive:    true,
				},
			},
			wantErr: false,
		},
	}
	for i := range tests {
		tt := tests[i]
		s.Run(tt.name, func() {
			got, err := s.testRepo.ListAllProducts(tt.args.ctx)
			if err != nil {
				if tt.wantErr == true {
					return
				}
				fmt.Printf("ListAllProducts() error = %v, wantErr %v", err, tt.wantErr)
				s.Fail("test ListAllProducts failed")
				return
			}
			if !s.Equal(tt.want, got) {
				fmt.Printf("ListAllProducts() got = %v, want %v", got, tt.want)
				s.Fail("test ListAllProducts failed")
			}
		})
	}
}

func (s *ProductTestSuite) Test_productRepo_SearchProductsByCategory() {
	type args struct {
		ctx      context.Context
		category int
	}
	tests := []struct {
		name    string
		args    args
		want    []model.Product
		wantErr bool
	}{
		{
			name: "search category with existing  product",
			args: args{
				ctx:      context.Background(),
				category: 1,
			},
			want: []model.Product{
				{
					ID:          s.productID,
					SKU:         "",
					Name:        "пушка",
					Type:        "",
					URI:         "/product/тепловая-3001",
					Description: "",
					IsActive:    false,
				},
			},
			wantErr: false,
		},
		{
			name: "search category with no  products",
			args: args{
				ctx:      context.Background(),
				category: 2,
			},
			want:    []model.Product{},
			wantErr: false,
		},
	}
	for i := range tests {
		tt := tests[i]
		s.Run(tt.name, func() {
			got, err := s.testRepo.SearchProductsByCategory(tt.args.ctx, tt.args.category)
			if err != nil {
				if tt.wantErr == true {
					return
				}
				fmt.Printf("SearchProductsByCategory() error = %v, wantErr %v", err, tt.wantErr)
				s.Fail("test SearchProductsByCategory failed")
				return
			}
			if !s.Equal(tt.want, got) {
				fmt.Printf("SearchProductsByCategory() got = %v, want %v", got, tt.want)
				s.Fail("test SearchProductsByCategory failed")
			}
		})
	}
}

func (s *ProductTestSuite) Test_productRepo_SearchProductsByName() {
	type args struct {
		ctx         context.Context
		productName string
	}
	tests := []struct {
		name    string
		args    args
		want    model.Product
		wantErr bool
	}{
		{
			name: "search existing product",
			args: args{
				ctx:         context.Background(),
				productName: "пушка",
			},
			want: model.Product{
				ID:          s.productID,
				SKU:         "3001",
				Name:        "пушка",
				Type:        "",
				URI:         "/product/тепловая-3001",
				Description: "пушка детская",
				IsActive:    false,
			},
			wantErr: false,
		},
		{
			name: "search non-existing product",
			args: args{
				ctx:         context.Background(),
				productName: "клюшка",
			},
			want: model.Product{
				ID:          "",
				SKU:         "",
				Name:        "",
				Type:        "",
				URI:         "",
				Description: "",
				IsActive:    false,
			},
			wantErr: false,
		},
	}
	for i := range tests {
		tt := tests[i]
		s.Run(tt.name, func() {
			got, err := s.testRepo.SearchProductsByName(tt.args.ctx, tt.args.productName)
			if err != nil {
				if tt.wantErr == true {
					return
				}
				fmt.Printf("SearchProductsByName() error = %v, wantErr %v", err, tt.wantErr)
				s.Fail("test SearchProductsByName failed")
				return
			}
			if !s.Equal(tt.want, got) {
				fmt.Printf("SearchProductsByName() got = %v, want %v", got, tt.want)
				s.Fail("test SearchProductsByName failed")
			}
		})
	}
}

func (s *ProductTestSuite) Test_productRepo_SearchProductsByShop() {
	type args struct {
		ctx    context.Context
		shopID int
	}
	tests := []struct {
		name    string
		args    args
		want    []model.Product
		wantErr bool
	}{
		{
			name: "search existing shop",
			args: args{
				ctx:    context.Background(),
				shopID: 1,
			},
			want: []model.Product{
				{
					ID:          s.productID,
					SKU:         "3001",
					Name:        "пушка",
					Type:        "",
					URI:         "/product/тепловая-3001",
					Description: "пушка детская",
					IsActive:    true,
				},
			},
			wantErr: false,
		},
		{
			name: "search non-existing shop",
			args: args{
				ctx:    context.Background(),
				shopID: 10,
			},
			want:    []model.Product{},
			wantErr: false,
		},
	}
	for i := range tests {
		tt := tests[i]
		s.Run(tt.name, func() {
			got, err := s.testRepo.SearchProductsByShop(tt.args.ctx, tt.args.shopID)
			if err != nil {
				if tt.wantErr == true {
					return
				}
				fmt.Printf("SearchProductsByShop() error = %v, wantErr %v", err, tt.wantErr)
				s.Fail("test SearchProductsByShop failed")
				return
			}
			if !s.Equal(tt.want, got) {
				fmt.Printf("SearchProductsByShop() got = %v, want %v", got, tt.want)
				s.Fail("test SearchProductsByShop failed")
			}
		})
	}
}
