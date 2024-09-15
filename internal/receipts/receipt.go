package receipts

import (
	log "github.com/sirupsen/logrus"
)

type DB interface {
	GetPoints(id string) (Points, error)
	Create(r Receipt, p Points) (Receipt, error)
}

type Service interface {
	GetPoints(id string) (Points, error)
	Create(receipt Receipt) (Receipt, error)
}

type receipt struct {
	db DB
}

func NewReceiptService(db DB) Service {
	return &receipt{
		db: db}
}

func (r *receipt) GetPoints(id string) (Points, error) {
	points, err := r.db.GetPoints(id)
	if err != nil {
		log.WithFields(log.Fields{
			"ID": id,
		}).Error("Failed to retrieve points for receipt")
		return Points{}, err
	}
	return points, nil
}

func (r *receipt) Create(receipt Receipt) (Receipt, error) {
	pointsObj := toPoints(receipt)

	createdReceipt, err := r.db.Create(receipt, pointsObj)
	if err != nil {
		return Receipt{}, ErrReceiptInvalid
	}

	return createdReceipt, nil
}
