package db

import (
	"fetch_take_home/internal/receipts"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDBGet(t *testing.T) {
	receipt := receipts.Receipt{
		ID:           "",
		Retailer:     "retailer",
		PurchaseDate: time.Now(),
		PurchaseTime: time.Now(),
		Items:        nil,
		Total:        0,
	}

	points := receipts.Points{
		ID:     "",
		Points: 0,
	}
	db := NewDB()
	createdReceipt, err := db.Create(receipt, points)
	assert.NoError(t, err)

	tests := map[string]struct {
		input  string
		expect receipts.Points
		err    error
	}{
		"Successful Get": {
			input:  createdReceipt.ID,
			expect: receipts.Points{},
			err:    nil,
		},
		"Not found error": {
			input:  "invalid",
			expect: receipts.Points{},
			err:    receipts.ErrReceiptNotFound,
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			response, err := db.GetPoints(test.input)

			assert.Equal(t, test.expect, response)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestDBCreate(t *testing.T) {
	purchaseDate, _ := time.Parse("2006-01-02", "2022-01-01")
	purchaseTime, _ := time.Parse("15:04", "13:01")
	receipt := receipts.Receipt{
		ID:           "",
		Retailer:     "retailer",
		PurchaseDate: purchaseDate,
		PurchaseTime: purchaseTime,
		Items:        nil,
		Total:        0,
	}

	points := receipts.Points{
		ID:     "",
		Points: 0,
	}

	db := NewDB()
	createdReceipt, err := db.Create(receipt, points)
	assert.NoError(t, err)
	assert.NotEqual(t, "", createdReceipt.ID)

	receipt.ID = createdReceipt.ID
	assert.Equal(t, receipt, createdReceipt)

	createdPoints, err := db.GetPoints(receipt.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, "", createdPoints.ID)

	points.ID = createdPoints.ID
	assert.Equal(t, points, createdPoints)
}
