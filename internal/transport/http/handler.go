package http

import (
	"fetch_take_home/errors"
	"fetch_take_home/internal/receipts"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	ReceiptService receipts.Service
}

func Activate(router *gin.Engine, receiptService receipts.Service) {
	handler := Handler{
		ReceiptService: receiptService,
	}

	router.GET("/receipts/:id/points", handler.GetPoints)
	router.POST("/receipts/process", handler.Create)
	router.GET("/health", handler.HealthCheck)
}

func getPointsResponse(p receipts.Points) gin.H {
	return gin.H{"points": p.Points}
}

func (h *Handler) GetPoints(c *gin.Context) {
	points, err := h.ReceiptService.GetPoints(c.Param("id"))
	if err != nil {
		status, e := handleError(err)
		c.IndentedJSON(status, e)
		return
	}
	c.IndentedJSON(http.StatusOK, getPointsResponse(points))
}

func createResponse(r receipts.Receipt) gin.H {
	return gin.H{"id": r.ID}
}

func (h *Handler) Create(c *gin.Context) {
	var receiptDTO receipts.ReceiptDTO

	if err := c.ShouldBindJSON(&receiptDTO); err != nil {
		status, e := handleError(receipts.ErrReceiptCreate)
		c.IndentedJSON(status, e)
		return
	}

	receipt, err := h.ReceiptService.Create(receiptDTO)
	if err != nil {
		status, e := handleError(err)
		c.IndentedJSON(status, e)
		return
	}

	c.IndentedJSON(http.StatusCreated, createResponse(receipt))
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "200", "healthy": "OK"})
}

func handleError(e error) (int, error) {
	switch e {
	case receipts.ErrReceiptNotFound:
		return http.StatusNotFound, errors.NewError(errors.NotFound, "No receipt found for that id")
	case receipts.ErrReceiptCreate:
		return http.StatusBadRequest, errors.NewError(errors.BadRequest, "The receipt is invalid")
	default:
		return http.StatusInternalServerError, errors.NewError(errors.InternalServerError, "Internal server error")
	}
}
