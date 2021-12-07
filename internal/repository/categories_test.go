package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
	"market4/internal/model"
	"testing"
)

type CategoriesTestSuite struct {
	suite.Suite
	testRepo categoryRepo
}

func Test_CategoriesSuite(t *testing.T) {
	suite.Run(t, new(CategoriesTestSuite))
}

func (suite *CategoriesTestSuite) SetupTest() {
	fmt.Println("start setup")
	var err error
	suite.testRepo.pool, err = pgxpool.Connect(context.Background(), testDSN)
	if err != nil {
		suite.Error(err)
		suite.Fail("setup failed")
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
		suite.Error(err)
		suite.Fail("setup failed")
		return
	}

	addCategoriesReq := "INSERT " +
		"INTO categories (name, uri_name) " +
		"VALUES ('Стройматериалы', 'Стройматериалы-1'), " +
		"('Игрушки', 'Игрушки-2');"

	_, err = suite.testRepo.pool.Exec(context.Background(), addCategoriesReq)
	if err != nil {
		suite.Error(err)
		suite.Fail("setup failed")
		return
	}
}

func (suite *CategoriesTestSuite) TearDownTest() {
	fmt.Println("cleaning up")
	var err error
	_, err = suite.testRepo.pool.Exec(context.Background(), "DROP TABLE categories CASCADE;")
	if err != nil {
		suite.Error(err)
		suite.Fail("cleaning failed")
	}
}

func (suite *CategoriesTestSuite) Test_categoryRepo_IfCategoryExists() {
	type args struct {
		ctx      context.Context
		category int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "existing category",
			args: args{
				ctx:      context.Background(),
				category: 1,
			},
			want: true,
		},
		{
			name: "non-existing category",
			args: args{
				ctx:      context.Background(),
				category: 10,
			},
			want: false,
		},
	}
	for i := range tests {
		tt := tests[i]
		suite.Run(tt.name, func() {
			got := suite.testRepo.IfCategoryExists(tt.args.ctx, tt.args.category)
			if got != tt.want {
				fmt.Printf("IfCategoryExists() = %v, want %v", got, tt.want)
				suite.Fail("test failed")
			}
		})
	}
}

func (suite *CategoriesTestSuite) Test_categoryRepo_ListAllCategories() {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []model.Category
		wantErr bool
	}{
		{
			name: "list of categories",
			args: args{
				ctx: context.Background(),
			},
			want: []model.Category{
				{
					ID:       1,
					Name:     "Стройматериалы",
					URI_name: "Стройматериалы-1",
				},
				{
					ID:       2,
					Name:     "Игрушки",
					URI_name: "Игрушки-2",
				},
			},
			wantErr: false,
		},
	}
	for i := range tests {
		tt := tests[i]
		suite.Run(tt.name, func() {
			got, err := suite.testRepo.ListAllCategories(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				fmt.Printf("ListAllCategories() error = %v, wantErr %v", err, tt.wantErr)
				suite.Fail("test ListAllCategories failed")
				return
			}
			if !suite.Equal(tt.want, got) {
				fmt.Printf("ListAllCategories() got = %v, want %v", got, tt.want)
				suite.Fail("test ListAllCategories failed")
			}
		})
	}
}

func (suite *CategoriesTestSuite) Test_categoryRepo_AddCategory() {
	type args struct {
		ctx      context.Context
		category *model.Category
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "add new category",
			args: args{
				ctx: context.Background(),
				category: &model.Category{
					Name: "Шуршики",
				},
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "add existing category",
			args: args{
				ctx: context.Background(),
				category: &model.Category{
					Name: "Игрушки",
				},
			},
			want:    0,
			wantErr: true,
		},
	}
	for i := range tests {
		tt := tests[i]
		suite.Run(tt.name, func() {
			got, err := suite.testRepo.AddCategory(tt.args.ctx, tt.args.category)
			if (err != nil) != tt.wantErr {
				fmt.Printf("AddCategory() error = %v, wantErr %v", err, tt.wantErr)
				suite.Fail("test AddCategory failed")
				return
			}
			if got != tt.want {
				fmt.Printf("AddCategory() got = %v, want %v", got, tt.want)
				suite.Fail("test AddCategory failed")
			}
		})
	}
}

func (suite *CategoriesTestSuite) Test_categoryRepo_EditCategory() {
	type args struct {
		ctx      context.Context
		category *model.Category
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "existing category",
			args: args{
				ctx: context.Background(),
				category: &model.Category{
					ID:   1,
					Name: "Шуршики",
				},
			},
			wantErr: false,
		},
		{
			name: "non-existing category",
			args: args{
				ctx: context.Background(),
				category: &model.Category{
					ID:   10,
					Name: "Шуршики",
				},
			},
			wantErr: true,
		},
	}
	for i := range tests {
		tt := tests[i]
		suite.Run(tt.name, func() {
			err := suite.testRepo.EditCategory(tt.args.ctx, tt.args.category)
			fmt.Printf("GOT: %v", err)
			if err != nil {
				if tt.wantErr == false {
					fmt.Printf("EditCategory() error = %v, wantErr %v", err, tt.wantErr)
					suite.Fail("test EditCategory failed")
					return
				}
			}
		})
	}
}
