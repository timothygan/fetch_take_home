package http

import (
	"fetch_take_home/internal/receipts"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func toItem(itemDTO receipts.ItemDTO) (receipts.Item, error) {
	val, err := strconv.ParseFloat(itemDTO.Price, 64)
	if err != nil {
		log.WithFields(log.Fields{
			"shortDescription": itemDTO.ShortDescription,
			"price":            itemDTO.Price,
		}).Error("Failed to parse item")
		return receipts.Item{}, receipts.ErrReceiptInvalid
	}
	cents := int64(val*100 + 0.5)

	return receipts.Item{
		ShortDescription: itemDTO.ShortDescription,
		Price:            cents,
	}, nil
}

func toReceipt(receiptDTO receipts.ReceiptDTO) (receipts.Receipt, error) {
	var purchaseDate, purchaseDateError = time.Parse("2006-01-02", receiptDTO.PurchaseDate)
	if purchaseDateError != nil {
		log.WithFields(log.Fields{
			"purchaseDate": receiptDTO.PurchaseDate,
		}).Error("Failed to parse purchase date")
		return receipts.Receipt{}, receipts.ErrReceiptInvalid
	}

	var purchaseTime, purchaseTimeError = time.Parse("15:04", receiptDTO.PurchaseTime)
	if purchaseTimeError != nil {
		log.WithFields(log.Fields{
			"purchaseTime": receiptDTO.PurchaseTime,
		}).Error("Failed to parse purchase time")
		return receipts.Receipt{}, receipts.ErrReceiptInvalid
	}

	var newItems []receipts.Item
	for _, itemDTO := range receiptDTO.Items {
		item, itemErr := toItem(itemDTO)
		if itemErr != nil {
			return receipts.Receipt{}, itemErr
		}
		newItems = append(newItems, item)
	}

	val, err := strconv.ParseFloat(receiptDTO.Total, 64)
	if err != nil {
		log.WithFields(log.Fields{
			"total": receiptDTO.Total,
		}).Error("Failed to parse total")
		return receipts.Receipt{}, receipts.ErrReceiptInvalid
	}
	cents := int64(val*100 + 0.5)

	return receipts.Receipt{
		ID:           "",
		Retailer:     receiptDTO.Retailer,
		PurchaseDate: purchaseDate,
		PurchaseTime: purchaseTime,
		Items:        newItems,
		Total:        cents,
	}, nil
}
