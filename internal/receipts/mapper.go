package receipts

import (
	log "github.com/sirupsen/logrus"
	"math"
	"regexp"
	"strings"
	"time"
)

var twoPM, _ = time.Parse("15:04", "14:00")
var fourPM, _ = time.Parse("15:04", "16:00")

func toPoints(receipt Receipt) Points {
	var points int64 = 0

	points += int64(len(regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(receipt.Retailer, "")))
	if receipt.Total%100 == 0 {
		points += 50
	}

	if receipt.Total%25 == 0 {
		points += 25
	}

	points += int64(len(receipt.Items) / 2 * 5)

	for _, item := range receipt.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			points += int64(math.Ceil(float64(item.Price) / 500))
			log.Info(int64(math.Ceil(float64(item.Price) / 500)))
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
