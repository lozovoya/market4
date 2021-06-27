package repository

import (
	"context"
	"fmt"
	"log"
	"market4/internal/model"
)

func (s *Storage) ListAllShops(ctx context.Context) ([]*model.Shop, error) {

	log.Println("list all shops controller repository")
	dbReq := "SELECT id, name, address, lon, lat, working_hours " +
		"FROM shops"

	shops := make([]*model.Shop, 0)
	rows, err := s.pool.Query(ctx, dbReq)
	if err != nil {
		log.Println(err)
		return shops, fmt.Errorf("ListAllShops: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var shop model.Shop
		err = rows.Scan(&shop.Id, &shop.Name, &shop.Address, &shop.Lon, &shop.Lat, &shop.WorkingHours)
		if err != nil {
			log.Println(err)
			return shops, fmt.Errorf("ListAllShops: %w", err)
		}
		shops = append(shops, &shop)
	}

	return shops, nil
}

func (m *Storage) AddShop(ctx context.Context, shop *model.Shop) (int, error) {
	dbReq := "INSERT " +
		"INTO shops (name, address, lon, lat, working_hours) " +
		"VALUES ($1, $2, $3, $4, $5) " +
		"RETURNING id"
	var id int
	err := m.pool.QueryRow(ctx,
		dbReq,
		shop.Name, shop.Address, shop.Lon, shop.Lat, shop.WorkingHours).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("AddShop: %w", err)
	}

	log.Printf("shop %s is added", shop.Name)
	return id, nil
}

func (s *Storage) EditShop(ctx context.Context, shop *model.Shop) error {
	var dbReq = "UPDATE shops SET "

	if !IsEmpty(shop.Name) {
		dbReq = fmt.Sprintf("%s name = '%s',", dbReq, shop.Name)
	}

	if !IsEmpty(shop.Address) {
		dbReq = fmt.Sprintf("%s address = '%s',", dbReq, shop.Address)
	}

	if !IsEmpty(shop.Lon) {
		dbReq = fmt.Sprintf("%s lon = '%s',", dbReq, shop.Lon)
	}

	if !IsEmpty(shop.Lat) {
		dbReq = fmt.Sprintf("%s lat = '%s',", dbReq, shop.Lat)
	}

	if !IsEmpty(shop.WorkingHours) {
		dbReq = fmt.Sprintf("%s working_hours = '%s',", dbReq, shop.WorkingHours)
	}

	dbReq = fmt.Sprintf("%s updated = CURRENT_TIMESTAMP WHERE id = %d", dbReq, shop.Id)
	_, err := s.pool.Exec(ctx, dbReq)
	if err != nil {
		return fmt.Errorf("UpdateShopParameter: %w", err)
	}

	log.Printf("shop %d is updated", shop.Id)
	return nil
}
