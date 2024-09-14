package receipts

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"math"
	"strconv"
	"strings"
	"time"
)

var twoPM, _ = time.Parse("15:04", "14:00")
var fourPM, _ = time.Parse("15:04", "16:00")

type DB interface {
	GetPoints(id string) (Points, error)
	Create(r Receipt, p Points) error
}

type Service interface {
	GetPoints(id string) (Points, error)
	Create(receiptDto ReceiptDTO) (Receipt, error)
}

type receipt struct {
	db DB
}

func NewReceiptService(db DB) Service {
	return &receipt{db}
}

func (r *receipt) GetPoints(id string) (Points, error) {
	points, err := r.db.GetPoints(id)
	if err != nil {
		log.WithFields(log.Fields{
			"ID": id,
		}).Error("Failed to retrieve points for receipt")
		return Points{}, errors.Join(ErrReceiptNotFound, err)
	}
	return points, nil
}

func (r *receipt) Create(receiptDto ReceiptDTO) (Receipt, error) {
	receiptObj, err := toReceipt(receiptDto)
	if err != nil {
		log.WithFields(log.Fields{
			"retailer":     receiptDto.Retailer,
			"purchaseDate": receiptDto.PurchaseDate,
			"purchaseTime": receiptDto.PurchaseTime,
			"items":        receiptDto.Items,
			"total":        receiptDto.Total,
		}).Error("Failed to create receipt")
		return Receipt{}, errors.Join(ErrReceiptCreate, err)
	}

	pointsObj := toPoints(receiptObj)

	err = r.db.Create(receiptObj, pointsObj)
	if err != nil {
		return Receipt{}, errors.Join(ErrReceiptCreate, err)
	}

	return receiptObj, nil
}

func toItem(itemDTO ItemDTO) (Item, error) {
	val, err := strconv.ParseFloat(itemDTO.Price, 64)
	if err != nil {
		log.WithFields(log.Fields{
			"shortDescription": itemDTO.ShortDescription,
			"price":            itemDTO.Price,
		}).Error("Failed to parse item")
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
		log.WithFields(log.Fields{
			"purchaseDate": receiptDTO.PurchaseDate,
		}).Error("Failed to parse purchase date")
		return Receipt{}, purchaseDateError
	}

	var purchaseTime, purchaseTimeError = time.Parse("15:04", receiptDTO.PurchaseTime)
	if purchaseTimeError != nil {
		log.WithFields(log.Fields{
			"purchaseTime": receiptDTO.PurchaseTime,
		}).Error("Failed to parse purchase time")
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
		log.WithFields(log.Fields{
			"total": receiptDTO.Total,
		}).Error("Failed to parse total")
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
			points += int64(math.Round(float64(item.Price) / 5000))
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
