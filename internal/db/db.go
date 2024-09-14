package db

import (
	"fetch_take_home/internal/receipts"
	"github.com/google/uuid"
)

type Database struct {
	pointsDB   map[string]*receipts.Points
	receiptsDB map[string]*receipts.Receipt
}

func NewDB() receipts.DB {
	pDB := make(map[string]*receipts.Points)
	rDB := make(map[string]*receipts.Receipt)

	return &Database{
		pointsDB:   pDB,
		receiptsDB: rDB,
	}
}

func (db *Database) GetPoints(id string) (receipts.Points, error) {
	if db.pointsDB[id] == nil {
		return receipts.Points{}, receipts.ErrReceiptNotFound
	}
	return *db.pointsDB[id], nil
}

func (db *Database) Create(r receipts.Receipt, p receipts.Points) (receipts.Receipt, error) {
	var id = uuid.NewString()
	db.receiptsDB[id] = &receipts.Receipt{
		ID:           id,
		Retailer:     r.Retailer,
		PurchaseDate: r.PurchaseDate,
		PurchaseTime: r.PurchaseDate,
		Items:        r.Items,
		Total:        r.Total,
	}
	db.pointsDB[id] = &p
	return *db.receiptsDB[id], nil
}
