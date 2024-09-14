package http

import (
	"errors"
	"fetch_take_home/internal/db"
	"fetch_take_home/internal/receipts"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	ReceiptService receipts.Service
}

func Setup(router *gin.Engine) {
	database := db.NewDB()
	receiptService := receipts.NewReceiptService(database)
	newHandler(router, receiptService)
}

func newHandler(router *gin.Engine, receiptService receipts.Service) {
	handler := Handler{
		ReceiptService: receiptService,
	}

	router.GET("/receipts/:id/points", handler.GetPoints)
	router.POST("/receipts/process", handler.Create)
}

func (h *Handler) GetPoints(c *gin.Context) {
	points, err := h.ReceiptService.GetPoints(c.Param("id"))
	if err != nil {
		status, e := handleError(err)
		c.IndentedJSON(status, e)
		return
	}
	c.IndentedJSON(http.StatusOK, points.Points)
}

func (h *Handler) Create(c *gin.Context) {
	var receiptDTO receipts.ReceiptDTO

	if err := c.ShouldBindJSON(&receiptDTO); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	receipt, err := h.ReceiptService.Create(receiptDTO)
	if err != nil {
		status, e := handleError(err)
		c.IndentedJSON(status, e)
		return
	}

	c.IndentedJSON(http.StatusCreated, receipt)
}

func handleError(e error) (int, error) {
	switch {
	case errors.Is(e, receipts.ErrReceiptNotFound):
		return http.StatusNotFound, e
	case errors.Is(e, receipts.ErrReceiptCreate):
		return http.StatusBadRequest, e
	default:
		return http.StatusInternalServerError, e
	}
}
