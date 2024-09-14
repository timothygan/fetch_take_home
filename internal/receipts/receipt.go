package receipts

import (
	"context"
	"errors"
	"math"
	"strconv"
	"strings"
	"time"
)

var twoPM, _ = time.Parse("15:04", "14:00")
var fourPM, _ = time.Parse("15:04", "16:00")

type DB interface {
	GetPoints(ctx context.Context, id string) (int64, error)
	Create(ctx context.Context, r Receipt, p Points) error
}

type Service interface {
	GetPoints(ctx context.Context, id string) (int64, error)
	Create(ctx context.Context, receiptDto ReceiptDTO) (string, error)
}

type receipt struct {
	db DB
}

func NewReceiptService(db DB) Service {
	return &receipt{db}
}

func (r *receipt) GetPoints(ctx context.Context, id string) (int64, error) {
	points, err := r.db.GetPoints(ctx, id)
	if err != nil {
		return -1, errors.Join(ErrReceiptNotFound, err)
	}
	return points, nil
}

func (r *receipt) Create(ctx context.Context, receiptDto ReceiptDTO) (string, error) {
	receiptObj, err := toReceipt(receiptDto)
	if err != nil {
		return "", errors.Join(ErrReceiptCreate, err)
	}

	pointsObj := toPoints(receiptObj)

	err = r.db.Create(ctx, receiptObj, pointsObj)
	if err != nil {
		return "", errors.Join(ErrReceiptCreate, err)
	}

	return receiptObj.ID, nil
}

func toItem(itemDTO ItemDTO) (Item, error) {
	val, err := strconv.ParseFloat(itemDTO.Price, 64)
	if err != nil {
		return Item{}, err
	}
	cents := int64(val*100 + 0.5)

	return Item{
		ShortDescription: itemDTO.ShortDescription,
		Price:            cents,
	}, nil
}

func toReceipt(receiptDTO ReceiptDTO) (Receipt, error) {
	var purchaseDate, purchaseDateError = time.Parse("2006-01-02", receiptDTO.PurchaseDate)
	if purchaseDateError != nil {
		return Receipt{}, purchaseDateError
	}

	var purchaseTime, purchaseTimeError = time.Parse("15:04", receiptDTO.PurchaseTime)
	if purchaseTimeError != nil {
		return Receipt{}, purchaseTimeError
	}

	var newItems []Item
	for _, itemDTO := range receiptDTO.Items {
		item, itemErr := toItem(itemDTO)
		if itemErr != nil {
			return Receipt{}, itemErr
		}
		newItems = append(newItems, item)
	}

	val, err := strconv.ParseFloat(receiptDTO.Total, 64)
	if err != nil {
		return Receipt{}, err
	}
	cents := int64(val*100 + 0.5)

	return Receipt{
		ID:           "",
		Retailer:     receiptDTO.Retailer,
		PurchaseDate: purchaseDate,
		PurchaseTime: purchaseTime,
		Items:        newItems,
		Total:        cents,
	}, nil
}

func toPoints(receipt Receipt) Points {
	var points int64 = 0

	points += int64(len(receipt.Retailer))

	if receipt.Total%100 == 0 {
		points += 50
	}

	if receipt.Total%25 == 0 {
		points += 25
	}

	points += int64(len(receipt.Items) / 2)

	for _, item := range receipt.Items {
		if len(strings.Trim(item.ShortDescription, "\\w"))%3 == 0 {
			points += int64(math.Round(float64(item.Price) / 1000 * .2))
		}
	}

	if receipt.PurchaseDate.Day()%2 == 1 {
		points += 6
	}

	if receipt.PurchaseTime.Before(fourPM) && receipt.PurchaseTime.After(twoPM) {
		points += 10
	}

	return Points{
		ID:     "",
		Points: points,
	}
}
