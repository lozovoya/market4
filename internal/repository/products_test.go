package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
	"market4/internal/model"
	"testing"
)

type ProductTestSuite struct {
	suite.Suite
	testRepo  productRepo
	productID string
}

func Test_ProductSuite(t *testing.T) {
	suite.Run(t, new(ProductTestSuite))
}

func (suite *ProductTestSuite) SetupTest() {
	fmt.Println("start setup")
	var err error

	suite.testRepo.pool, err = pgxpool.Connect(context.Background(), testDSN)
	if err != nil {
		suite.Fail("setup failed", err)
		return
	}

	createExtensionReq := "CREATE EXTENSION pgcrypto;"
	_, err = suite.testRepo.pool.Exec(context.Background(), createExtensionReq)
	if err != nil {
		fmt.Println("pgcrypto failed: createExtensionReq", err)
	}
	createTableProductsReq := "CREATE " +
		"TABLE products ( " +
		"id          UUID DEFAULT gen_random_uuid() PRIMARY KEY, " +
		"sku         TEXT NOT NULL, " +
		"name        TEXT NOT NULL, " +
		"uri         TEXT NOT NULL, " +
		"description TEXT NOT NULL, " +
		"is_active       BOOL NOT NULL, " +
		"created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, " +
		"updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);"
	_, err = suite.testRepo.pool.Exec(context.Background(), createTableProductsReq)
	if err != nil {
		suite.Fail("setup failed: createTableProductsReq", err)
		return
	}

	addProductReq := "INSERT " +
		"INTO products (sku, name, uri, description, is_active) " +
		"VALUES ('3001', 'пушка', '/product/тепловая-3001', 'пушка детская', true) RETURNING id;"
	err = suite.testRepo.pool.QueryRow(context.Background(), addProductReq).Scan(&suite.productID)
	if err != nil {
		suite.Fail("setup failed: addProductReq", err)
		return
	}
	createTableCategoriesReq := "CREATE " +
		"TABLE categories ( " +
		"id BIGSERIAL PRIMARY KEY, " +
		"name TEXT NOT NULL UNIQUE, " +
		"uri_name TEXT UNIQUE, " +
		"created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, " +
		"updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);"
	_, err = suite.testRepo.pool.Exec(context.Background(), createTableCategoriesReq)
	if err != nil {
		suite.Fail("setup failed", err)
		return
	}

	addCategoriesReq := "INSERT " +
		"INTO categories (name, uri_name) " +
		"VALUES ('Стройматериалы', 'Стройматериалы-1'), " +
		"('Игрушки', 'Игрушки-2');"

	_, err = suite.testRepo.pool.Exec(context.Background(), addCategoriesReq)
	if err != nil {
		suite.Fail("setup failed", err)
		return
	}
	createTableProductCategoryReq := "CREATE " +
		"TABLE productcategory ( " +
		"category_id  BIGINT NOT NULL REFERENCES categories, " +
		"product_id UUID NOT NULL REFERENCES products, " +
		"PRIMARY KEY (category_id, product_id));"
	_, err = suite.testRepo.pool.Exec(context.Background(), createTableProductCategoryReq)
	if err != nil {
		suite.Fail("setup failed", err)
		return
	}
	addProductCategoryReq := fmt.Sprintf("INSERT INTO productcategory (category_id, product_id) VALUES (1, '%s');", suite.productID)
	_, err = suite.testRepo.pool.Exec(context.Background(), addProductCategoryReq)
	if err != nil {
		suite.Fail("setup failed", err)
		return
	}
	createTableShopsReq := "CREATE " +
		"TABLE shops ( " +
		"id              BIGSERIAL PRIMARY KEY, " +
		"name            TEXT NOT NULL, " +
		"address         TEXT NOT NULL, " +
		"lon             TEXT, " +
		"lat             TEXT, " +
		"working_hours   TEXT, " +
		"created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, " +
		"updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);"
	_, err = suite.testRepo.pool.Exec(context.Background(), createTableShopsReq)
	if err != nil {
		suite.Fail("setup failed", err)
		return
	}
	addShopReq := "INSERT " +
		"INTO shops (name, address, lon, lat, working_hours) " +
		"VALUES ('Магазин на диване', 'Москва, Останкино', '324234' , '5465476', '8 - 20'), " +
		"('Магазин для взрослых', 'Ростов, кремль', '12334' , '5465476', '8 - 20');"
	_, err = suite.testRepo.pool.Exec(context.Background(), addShopReq)
	if err != nil {
		suite.Fail("setup failed", err)
		return
	}

	createTableProductShopReq := "CREATE " +
		"TABLE productshop ( " +
		"shop_id BIGINT NOT NULL REFERENCES shops, " +
		"product_id UUID NOT NULL REFERENCES products, " +
		"PRIMARY KEY (shop_id, product_id));"
	_, err = suite.testRepo.pool.Exec(context.Background(), createTableProductShopReq)
	if err != nil {
		suite.Fail("setup failed", err)
		return
	}
	addProductShopReq := fmt.Sprintf("INSERT INTO productshop (shop_id, product_id) VALUES (1, '%s');", suite.productID)
	_, err = suite.testRepo.pool.Exec(context.Background(), addProductShopReq)
	if err != nil {
		suite.Fail("setup failed", err)
		return
	}
}

func (suite *ProductTestSuite) TearDownTest() {
	fmt.Println("cleaning up")
	var err error
	_, err = suite.testRepo.pool.Exec(context.Background(),
		"DROP TABLE products, categories, shops, productshop, productcategory CASCADE;")
	if err != nil {
		suite.Error(err)
		suite.Fail("cleaning failed")
	}
}

func (suite *ProductTestSuite) Test_productRepo_IfProductExists() {
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
				productID: suite.productID,
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
		suite.Run(tt.name, func() {
			if got := suite.testRepo.IfProductExists(tt.args.ctx, tt.args.productID); got != tt.want {
				fmt.Printf("IfProductExists() = %v, want %v", got, tt.want)
				suite.Fail("test IfProductExists failed")
			}
		})
	}
}

func (suite *ProductTestSuite) Test_productRepo_setProductCategory() {
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
				productId:  suite.productID,
			},
			wantErr: false,
		},
		{
			name: "set non-existing category",
			args: args{
				ctx:        context.Background(),
				categoryId: 20,
				productId:  suite.productID,
			},
			wantErr: true,
		},
	}
	for i := range tests {
		tt := tests[i]
		suite.Run(tt.name, func() {
			err := suite.testRepo.setProductCategory(tt.args.ctx, tt.args.categoryId, tt.args.productId)
			if err != nil {
				if tt.wantErr == true {
					return
				}
				fmt.Printf("setProductCategory() error = %v, wantErr %v", err, tt.wantErr)
				suite.Fail("test setProductCategory failed")
			}
		})
	}
}

func (suite *ProductTestSuite) Test_productRepo_setProductShop() {
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
				productID: suite.productID,
			},
			wantErr: false,
		},
		{
			name: "set non-existing shop",
			args: args{
				ctx:       context.Background(),
				shopID:    10,
				productID: suite.productID,
			},
			wantErr: true,
		},
	}
	for i := range tests {
		tt := tests[i]
		suite.Run(tt.name, func() {
			err := suite.testRepo.setProductShop(tt.args.ctx, tt.args.shopID, tt.args.productID)
			if err != nil {
				if tt.wantErr == true {
					return
				}
				fmt.Printf("setProductShop() error = %v, wantErr %v", err, tt.wantErr)
				suite.Fail("test setProductShop failed")
			}
		})
	}
}

func (suite *ProductTestSuite) Test_productRepo_EditProduct() {
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
				ID:   suite.productID,
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
		suite.Run(tt.name, func() {
			got, err := suite.testRepo.EditProduct(tt.args.ctx, tt.args.product, tt.args.shopID, tt.args.categoryID)
			if err != nil {
				if tt.wantErr == true {
					return
				}
				fmt.Printf("EditProduct() error = %v, wantErr %v", err, tt.wantErr)
				suite.Fail("test EditProduct failed")
				return
			}
			if !suite.Equal(tt.want, got) {
				fmt.Printf("EditProduct() got = %v, want %v", got, tt.want)
				suite.Fail("test EditProduct failed")
			}
		})
	}
}

func (suite *ProductTestSuite) Test_productRepo_ListAllProducts() {
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
					ID:          suite.productID,
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
		suite.Run(tt.name, func() {
			got, err := suite.testRepo.ListAllProducts(tt.args.ctx)
			if err != nil {
				if tt.wantErr == true {
					return
				}
				fmt.Printf("ListAllProducts() error = %v, wantErr %v", err, tt.wantErr)
				suite.Fail("test ListAllProducts failed")
				return
			}
			if !suite.Equal(tt.want, got) {
				fmt.Printf("ListAllProducts() got = %v, want %v", got, tt.want)
				suite.Fail("test ListAllProducts failed")
			}
		})
	}
}

func (suite *ProductTestSuite) Test_productRepo_SearchProductsByCategory() {
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
					ID:          suite.productID,
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
		suite.Run(tt.name, func() {
			got, err := suite.testRepo.SearchProductsByCategory(tt.args.ctx, tt.args.category)
			if err != nil {
				if tt.wantErr == true {
					return
				}
				fmt.Printf("SearchProductsByCategory() error = %v, wantErr %v", err, tt.wantErr)
				suite.Fail("test SearchProductsByCategory failed")
				return
			}
			if !suite.Equal(tt.want, got) {
				fmt.Printf("SearchProductsByCategory() got = %v, want %v", got, tt.want)
				suite.Fail("test SearchProductsByCategory failed")
			}
		})
	}
}

func (suite *ProductTestSuite) Test_productRepo_SearchProductsByName() {
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
				ID:          suite.productID,
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
		suite.Run(tt.name, func() {
			got, err := suite.testRepo.SearchProductsByName(tt.args.ctx, tt.args.productName)
			if err != nil {
				if tt.wantErr == true {
					return
				}
				fmt.Printf("SearchProductsByName() error = %v, wantErr %v", err, tt.wantErr)
				suite.Fail("test SearchProductsByName failed")
				return
			}
			if !suite.Equal(tt.want, got) {
				fmt.Printf("SearchProductsByName() got = %v, want %v", got, tt.want)
				suite.Fail("test SearchProductsByName failed")
			}
		})
	}
}

func (suite *ProductTestSuite) Test_productRepo_SearchProductsByShop() {
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
					ID:          suite.productID,
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
		suite.Run(tt.name, func() {
			got, err := suite.testRepo.SearchProductsByShop(tt.args.ctx, tt.args.shopID)
			if err != nil {
				if tt.wantErr == true {
					return
				}
				fmt.Printf("SearchProductsByShop() error = %v, wantErr %v", err, tt.wantErr)
				suite.Fail("test SearchProductsByShop failed")
				return
			}
			if !suite.Equal(tt.want, got) {
				fmt.Printf("SearchProductsByShop() got = %v, want %v", got, tt.want)
				suite.Fail("test SearchProductsByShop failed")
			}
		})
	}
}
