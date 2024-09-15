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
	validItemsDTO := []ItemDTO{
		{
			ShortDescription: "chicken",
			Price:            "5.00",
		},
	}
	validItems := []Item{
		{
			ShortDescription: "chicken",
			Price:            500,
		},
	}
	invalidItemsDTO := []ItemDTO{
		{
			ShortDescription: "chicken",
			Price:            "5.00",
		},
	}

	validReceiptDTO := ReceiptDTO{
		Retailer:     "retailer",
		PurchaseDate: "2024-09-14",
		PurchaseTime: "04:00",
		Items:        validItemsDTO,
		Total:        "5.00",
	}

	invalidReceiptDateDTO := ReceiptDTO{
		Retailer:     "retailer",
		PurchaseDate: "invalid date",
		PurchaseTime: "04:00",
		Items:        validItemsDTO,
		Total:        "5.00",
	}

	invalidReceiptTimeDTO := ReceiptDTO{
		Retailer:     "retailer",
		PurchaseDate: "2024-09-14",
		PurchaseTime: "invalid time",
		Items:        validItemsDTO,
		Total:        "5.00",
	}

	invalidReceiptItemsDTO := ReceiptDTO{
		Retailer:     "retailer",
		PurchaseDate: "2024-09-14",
		PurchaseTime: "04:00",
		Items:        invalidItemsDTO,
		Total:        "5.00",
	}

	invalidReceiptTotalDTO := ReceiptDTO{
		Retailer:     "retailer",
		PurchaseDate: "2024-09-14",
		PurchaseTime: "04:00",
		Items:        validItemsDTO,
		Total:        "invalid total",
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
		input  ReceiptDTO
		err    error
	}{
		"Success": {
			db: &dbMock{
				CreateResult: createdReceipt,
				CreateError:  nil,
			},
			input:  validReceiptDTO,
			result: createdReceipt,
			err:    nil,
		},
		"Invalid Receipt (invalid purchase date)": {
			db: &dbMock{
				CreateResult: Receipt{},
				CreateError:  ErrReceiptInvalid,
			},
			input:  invalidReceiptDateDTO,
			result: Receipt{},
			err:    ErrReceiptInvalid,
		},
		"Invalid Receipt (invalid purchase time)": {
			db: &dbMock{
				CreateResult: Receipt{},
				CreateError:  ErrReceiptInvalid,
			},
			input:  invalidReceiptTimeDTO,
			result: Receipt{},
			err:    ErrReceiptInvalid,
		},
		"Invalid Receipt (invalid purchase items)": {
			db: &dbMock{
				CreateResult: Receipt{},
				CreateError:  ErrReceiptInvalid,
			},
			input:  invalidReceiptItemsDTO,
			result: Receipt{},
			err:    ErrReceiptInvalid,
		},
		"Invalid Receipt (invalid purchase total)": {
			db: &dbMock{
				CreateResult: Receipt{},
				CreateError:  ErrReceiptInvalid,
			},
			input:  invalidReceiptTotalDTO,
			result: Receipt{},
			err:    ErrReceiptInvalid,
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

func TestReceiptToReceipt(t *testing.T) {
	targetPurchaseDate, _ := time.Parse("2006-01-02", "2022-01-01")
	targetPurchaseTime, _ := time.Parse("15:04", "13:01")
	cornerMarketPurchaseDate, _ := time.Parse("2006-01-02", "2022-03-20")
	cornerMarketPurchaseTime, _ := time.Parse("15:04", "14:33")
	targetItemsDTO := []ItemDTO{
		{
			ShortDescription: "Mountain Dew 12PK",
			Price:            "6.49",
		}, {
			ShortDescription: "Emils Cheese Pizza",
			Price:            "12.25",
		}, {
			ShortDescription: "Knorr Creamy Chicken",
			Price:            "1.26",
		}, {
			ShortDescription: "Doritos Nacho Cheese",
			Price:            "3.35",
		}, {
			ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
			Price:            "12.00",
		},
	}
	cornerMarketItemsDTO := []ItemDTO{
		{
			ShortDescription: "Gatorade",
			Price:            "2.25",
		}, {
			ShortDescription: "Gatorade",
			Price:            "2.25",
		}, {
			ShortDescription: "Gatorade",
			Price:            "2.25",
		}, {
			ShortDescription: "Gatorade",
			Price:            "2.25",
		},
	}
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

	targetReceiptDTO := ReceiptDTO{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items:        targetItemsDTO,
		Total:        "35.35",
	}

	cornerMarketDTO := ReceiptDTO{
		Retailer:     "M&M Corner Market",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Items:        cornerMarketItemsDTO,
		Total:        "9.00",
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

	tests := map[string]struct {
		result Receipt
		input  ReceiptDTO
	}{
		"Target receipt mapped successfully": {
			input:  targetReceiptDTO,
			result: targetReceipt,
		},

		"Corner Market receipt  mapped successfully": {
			input:  cornerMarketDTO,
			result: cornerMarketReceipt,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			response, _ := toReceipt(test.input)

			assert.Equal(t, test.result, response)
		})
	}
}
