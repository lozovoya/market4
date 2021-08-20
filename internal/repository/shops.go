package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"market4/internal/model"
	"strings"
)

type shopRepo struct {
	pool *pgxpool.Pool
}

func NewShopRepository(pool *pgxpool.Pool) Shop {
	return &shopRepo{pool: pool}
}

func (s *shopRepo) IfShopExists(ctx context.Context, shop int) bool {

	dbReq := "SELECT id FROM shops WHERE id=$1"
	var id = 0
	err := s.pool.QueryRow(ctx, dbReq, shop).Scan(&id)
	if err != nil {
		log.Println(fmt.Errorf("ifShopExists: %w", err))
		return false
	}
	if id != 0 {
		return true
	}
	return false
}

func (s *shopRepo) ListAllShops(ctx context.Context) ([]*model.Shop, error) {

	dbReq := "SELECT id, name, address, lon, lat, working_hours " +
		"FROM shops"

	shops := make([]*model.Shop, 0)
	rows, err := s.pool.Query(ctx, dbReq)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return shops, nil
		}
		return shops, fmt.Errorf("ListAllShops: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var shop model.Shop
		err = rows.Scan(&shop.ID, &shop.Name, &shop.Address, &shop.LON, &shop.LAT, &shop.WorkingHours)
		if err != nil {
			log.Println(err)
			return shops, fmt.Errorf("ListAllShops: %w", err)
		}
		shops = append(shops, &shop)
	}

	return shops, nil
}

func (s *shopRepo) AddShop(ctx context.Context, shop *model.Shop) (int, error) {
	dbReq := "INSERT " +
		"INTO shops (name, address, lon, lat, working_hours) " +
		"VALUES ($1, $2, $3, $4, $5) " +
		"RETURNING id"
	var id int
	err := s.pool.QueryRow(ctx,
		dbReq,
		shop.Name, shop.Address, shop.LON, shop.LAT, shop.WorkingHours).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("AddShop: %w", err)
	}

	log.Printf("shop %s is added", shop.Name)
	return id, nil
}

func (s *shopRepo) EditShop(ctx context.Context, shop *model.Shop) error {
	var dbReq = "UPDATE shops SET "

	if !IsEmpty(shop.Name) {
		dbReq = fmt.Sprintf("%s name = '%s',", dbReq, shop.Name)
	}

	if !IsEmpty(shop.Address) {
		dbReq = fmt.Sprintf("%s address = '%s',", dbReq, shop.Address)
	}

	if !IsEmpty(shop.LON) {
		dbReq = fmt.Sprintf("%s lon = '%s',", dbReq, shop.LON)
	}

	if !IsEmpty(shop.LAT) {
		dbReq = fmt.Sprintf("%s lat = '%s',", dbReq, shop.LAT)
	}

	if !IsEmpty(shop.WorkingHours) {
		dbReq = fmt.Sprintf("%s working_hours = '%s',", dbReq, shop.WorkingHours)
	}

	dbReq = fmt.Sprintf("%s updated = CURRENT_TIMESTAMP WHERE id = %d", dbReq, shop.ID)
	_, err := s.pool.Exec(ctx, dbReq)
	if err != nil {
		return fmt.Errorf("UpdateShopParameter: %w", err)
	}

	log.Printf("shop %d is updated", shop.ID)
	return nil
}
