package db

import (
	"context"
	"fetch_take_home/internal/receipts"
	"github.com/google/uuid"
)

type Database struct {
	pointsDB   map[string]*int64
	receiptsDB map[string]*receipts.Receipt
}

func NewDB() receipts.DB {
	pDB := make(map[string]*int64)
	rDB := make(map[string]*receipts.Receipt)

	return &Database{
		pointsDB:   pDB,
		receiptsDB: rDB,
	}
}

func (db *Database) GetPoints(ctx context.Context, id string) (int64, error) {
	return *db.pointsDB[id], nil
}

func (db *Database) Create(ctx context.Context, r receipts.Receipt, p receipts.Points) error {
	var id = uuid.NewString()
	db.receiptsDB[id] = &receipts.Receipt{
		ID:           id,
		Retailer:     r.Retailer,
		PurchaseDate: r.PurchaseDate,
		PurchaseTime: r.PurchaseDate,
		Items:        r.Items,
		Total:        r.Total,
	}
	db.pointsDB[id] = &p.Points
	return nil
}
