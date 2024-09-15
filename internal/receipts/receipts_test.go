package receipts

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type dbMock struct {
	GetPointsResult Points
	GetError        error

	CreateResult Receipt
	CreateError  error

	CreatePoints Points
}

func (db *dbMock) GetPoints(id string) (Points, error) {
	return db.GetPointsResult, db.GetError
}

func (db *dbMock) Create(r Receipt, p Points) (Receipt, error) {
	return db.CreateResult, db.CreateError
}

func TestReceiptServiceGetPoints(t *testing.T) {
	id := uuid.NewString()
	tests := map[string]struct {
		db     DB
		result Points
		err    error
	}{
		"Successfully retrieves Points": {
			db: &dbMock{
				GetPointsResult: Points{ID: id},
				GetError:        nil,
			},
			result: Points{ID: id},
			err:    nil,
		},
		"Receipt not found": {
			db: &dbMock{
				GetPointsResult: Points{},
				GetError:        ErrReceiptNotFound,
			},
			result: Points{},
			err:    ErrReceiptNotFound,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			service := NewReceiptService(test.db)
			response, err := service.GetPoints(id)

			assert.Equal(t, test.result, response)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestReceiptServiceCreate(t *testing.T) {
	id := uuid.NewString()
	purchaseDate, _ := time.Parse("2006-01-02", "2024-09-14")
	purchaseTime, _ := time.Parse("15:04", "14:00")

	validItems := []Item{
		{
			ShortDescription: "chicken",
			Price:            500,
		},
	}

	validReceipt := Receipt{
		ID:           "",
		Retailer:     "retailer",
		PurchaseDate: purchaseDate,
		PurchaseTime: purchaseTime,
		Items:        validItems,
		Total:        500,
	}

	createdReceipt := Receipt{
		ID:           id,
		Retailer:     "retailer",
		PurchaseDate: purchaseDate,
		PurchaseTime: purchaseTime,
		Items:        validItems,
		Total:        500,
	}

	tests := map[string]struct {
		db     DB
		result Receipt
		input  Receipt
		err    error
	}{
		"Success": {
			db: &dbMock{
				CreateResult: createdReceipt,
				CreateError:  nil,
			},
			input:  validReceipt,
			result: createdReceipt,
			err:    nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			service := NewReceiptService(test.db)
			response, err := service.Create(test.input)

			assert.Equal(t, test.result, response)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestReceiptToPoints(t *testing.T) {
	targetPurchaseDate, _ := time.Parse("2006-01-02", "2022-01-01")
	targetPurchaseTime, _ := time.Parse("15:04", "13:01")
	cornerMarketPurchaseDate, _ := time.Parse("2006-01-02", "2022-03-20")
	cornerMarketPurchaseTime, _ := time.Parse("15:04", "14:33")
	targetItems := []Item{
		{
			ShortDescription: "Mountain Dew 12PK",
			Price:            649,
		}, {
			ShortDescription: "Emils Cheese Pizza",
			Price:            1225,
		}, {
			ShortDescription: "Knorr Creamy Chicken",
			Price:            126,
		}, {
			ShortDescription: "Doritos Nacho Cheese",
			Price:            335,
		}, {
			ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
			Price:            1200,
		},
	}
	cornerMarketItems := []Item{
		{
			ShortDescription: "Gatorade",
			Price:            225,
		}, {
			ShortDescription: "Gatorade",
			Price:            225,
		}, {
			ShortDescription: "Gatorade",
			Price:            225,
		}, {
			ShortDescription: "Gatorade",
			Price:            225,
		},
	}

	targetReceipt := Receipt{
		Retailer:     "Target",
		PurchaseDate: targetPurchaseDate,
		PurchaseTime: targetPurchaseTime,
		Items:        targetItems,
		Total:        3535,
	}

	cornerMarketReceipt := Receipt{
		Retailer:     "M&M Corner Market",
		PurchaseDate: cornerMarketPurchaseDate,
		PurchaseTime: cornerMarketPurchaseTime,
		Items:        cornerMarketItems,
		Total:        900,
	}

	targetPoints := Points{
		ID:     "",
		Points: 28,
	}

	cornerMarketPoints := Points{
		ID:     "",
		Points: 109,
	}

	tests := map[string]struct {
		result Points
		input  Receipt
	}{
		"Target receipt points are correct": {
			input:  targetReceipt,
			result: targetPoints,
		},
		"Corner Market receipt points are correct": {
			input:  cornerMarketReceipt,
			result: cornerMarketPoints,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			points := toPoints(test.input)

			assert.Equal(t, test.result, points)
		})
	}
}
