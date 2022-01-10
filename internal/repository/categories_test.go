package repository

import (
	"context"
	"fmt"
	"market4/internal/model"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

type CategoriesTestSuite struct {
	suite.Suite
	testRepo categoryRepo
	Data     TestData
}

func Test_CategoriesSuite(t *testing.T) {
	suite.Run(t, new(CategoriesTestSuite))
}

func (s *CategoriesTestSuite) SetupTest() {
	fmt.Println("start setup")
	var err error
	s.testRepo.pool, err = pgxpool.Connect(context.Background(), testDSN)
	if err != nil {
		s.Error(err)
		s.Fail("setup failed")
		return
	}
	s.Data, err = loadTestDataFromYaml("categories_test.yaml")
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

func (s *CategoriesTestSuite) TearDownTest() {
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

func (s *CategoriesTestSuite) Test_categoryRepo_IfCategoryExists() {
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
		s.Run(tt.name, func() {
			got := s.testRepo.IfCategoryExists(tt.args.ctx, tt.args.category)
			if got != tt.want {
				fmt.Printf("IfCategoryExists() = %v, want %v", got, tt.want)
				s.Fail("test failed")
			}
		})
	}
}

func (s *CategoriesTestSuite) Test_categoryRepo_ListAllCategories() {
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
		s.Run(tt.name, func() {
			got, err := s.testRepo.ListAllCategories(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				fmt.Printf("ListAllCategories() error = %v, wantErr %v", err, tt.wantErr)
				s.Fail("test ListAllCategories failed")
				return
			}
			if !s.Equal(tt.want, got) {
				fmt.Printf("ListAllCategories() got = %v, want %v", got, tt.want)
				s.Fail("test ListAllCategories failed")
			}
		})
	}
}

func (s *CategoriesTestSuite) Test_categoryRepo_AddCategory() {
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
		s.Run(tt.name, func() {
			got, err := s.testRepo.AddCategory(tt.args.ctx, tt.args.category)
			if (err != nil) != tt.wantErr {
				fmt.Printf("AddCategory() error = %v, wantErr %v", err, tt.wantErr)
				s.Fail("test AddCategory failed")
				return
			}
			if got != tt.want {
				fmt.Printf("AddCategory() got = %v, want %v", got, tt.want)
				s.Fail("test AddCategory failed")
			}
		})
	}
}

func (s *CategoriesTestSuite) Test_categoryRepo_EditCategory() {
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
		s.Run(tt.name, func() {
			err := s.testRepo.EditCategory(tt.args.ctx, tt.args.category)
			fmt.Printf("GOT: %v", err)
			if err != nil {
				if tt.wantErr == false {
					fmt.Printf("EditCategory() error = %v, wantErr %v", err, tt.wantErr)
					s.Fail("test EditCategory failed")
					return
				}
			}
		})
	}
}
