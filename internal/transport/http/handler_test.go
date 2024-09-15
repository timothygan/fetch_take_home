package http

import (
	"encoding/json"
	"fetch_take_home/errors"
	"fetch_take_home/internal/receipts"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockReceiptService struct {
	GetPointsResult receipts.Points
	GetPointsError  error

	CreateResult receipts.Receipt
	CreateError  error
}

func (s *mockReceiptService) GetPoints(id string) (receipts.Points, error) {
	return s.GetPointsResult, s.GetPointsError
}

func (s *mockReceiptService) Create(receiptDTO receipts.ReceiptDTO) (receipts.Receipt, error) {
	return s.CreateResult, s.CreateError
}

func TestHandlerGetPoints(t *testing.T) {
	id := uuid.NewString()
	tests := map[string]struct {
		mockService receipts.Service
		uri         string
		response    interface{}
		statusCode  int
	}{
		"Successful Get": {
			mockService: &mockReceiptService{
				GetPointsResult: receipts.Points{ID: id, Points: 0},
				GetPointsError:  nil,
			},
			uri:        fmt.Sprintf("/receipts/%s/points", id),
			response:   receipts.PointsResponse{Points: 0},
			statusCode: http.StatusOK,
		},
		"ID not found": {
			mockService: &mockReceiptService{
				GetPointsResult: receipts.Points{},
				GetPointsError:  receipts.ErrReceiptNotFound,
			},
			uri: fmt.Sprintf("/receipts/%s/points", "invalid_id"),
			response: errors.AppError{
				Code:        "404",
				Description: "No receipt found for that id",
			},
			statusCode: http.StatusNotFound,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			response := httptest.NewRecorder()
			router := gin.New()
			Activate(router, test.mockService)

			req, err := http.NewRequest(http.MethodGet, test.uri, nil)
			assert.NoError(t, err)

			router.ServeHTTP(response, req)

			assert.Equal(t, test.statusCode, response.Code)

			if test.statusCode == http.StatusOK {
				var p receipts.PointsResponse
				if err := json.Unmarshal(response.Body.Bytes(), &p); err != nil {
					assert.Fail(t, "failed to unmarshal", response.Body.String(), err)
				}
				assert.Equal(t, test.response, p)
			} else {
				var err errors.AppError
				if err := json.Unmarshal(response.Body.Bytes(), &err); err != nil {
					assert.Fail(t, "failed to unmarshal", response.Body.String(), err)
				}
				assert.Equal(t, test.response, err)
			}
		})
	}
}

func TestHandlerCreate(t *testing.T) {
	id := uuid.NewString()
	retailer := "Walgreens"
	purchaseDate := "2022-01-02"
	purchaseTime := "08:13"
	total := "2.65"
	items := `[{"shortDescription": "Pepsi - 12-oz", "price": "1.25"},{"shortDescription": "Dasani", "price": "1.40"}]`
	receiptDTO := fmt.Sprintf(`{"retailer": "%s","purchaseDate": "%s","purchaseTime": "%s","total": "%s","items": %s}`,
		retailer,
		purchaseDate,
		purchaseTime,
		total,
		items,
	)

	tests := map[string]struct {
		mockService receipts.Service
		uri         string
		body        string
		response    interface{}
		statusCode  int
	}{
		"Successful Create": {
			mockService: &mockReceiptService{
				CreateResult: receipts.Receipt{ID: id},
				CreateError:  nil,
			},
			uri:        "/receipts/process",
			body:       receiptDTO,
			response:   receipts.CreateResponse{ID: id},
			statusCode: http.StatusOK,
		},
		"Invalid Receipt": {
			mockService: &mockReceiptService{
				CreateResult: receipts.Receipt{},
				CreateError:  receipts.ErrReceiptInvalid,
			},
			uri:  "/receipts/process",
			body: "{}",
			response: errors.AppError{
				Code:        "400",
				Description: "The receipt is invalid",
			},
			statusCode: http.StatusBadRequest,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			response := httptest.NewRecorder()
			router := gin.New()
			Activate(router, test.mockService)

			req, err := http.NewRequest(http.MethodPost, test.uri, strings.NewReader(test.body))
			assert.NoError(t, err)

			router.ServeHTTP(response, req)

			assert.Equal(t, test.statusCode, response.Code)
			if test.statusCode == http.StatusOK {
				var c receipts.CreateResponse
				if err := json.Unmarshal(response.Body.Bytes(), &c); err != nil {
					assert.Fail(t, "failed to unmarshal", response.Body.String(), err)
				}
				assert.Equal(t, test.response, c)
			} else {
				var err errors.AppError
				if err := json.Unmarshal(response.Body.Bytes(), &err); err != nil {
					assert.Fail(t, "failed to unmarshal", response.Body.String(), err)
				}
				assert.Equal(t, test.response, err)
			}
		})
	}
}
