package main

import (
	"fetch_take_home/internal/db"
	"fetch_take_home/internal/receipts"
	"fetch_take_home/internal/transport/http"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Run() error {
	database := db.NewDB()
	service := receipts.NewReceiptService(database)
	router := gin.New()
	http.Activate(router, service)
	if err := router.Run("localhost:8080"); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Error(err)
		log.Fatal("Failed to start server")
	}
}
