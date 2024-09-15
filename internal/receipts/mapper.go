package receipts

import (
	"math"
	"regexp"
	"strings"
	"time"
)

var twoPM, _ = time.Parse("15:04", "14:00")
var fourPM, _ = time.Parse("15:04", "16:00")
var alphanumeric = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func toPoints(receipt Receipt) Points {
	var points int64 = 0

	// Add points for alphanumeric characters
	points += int64(len(alphanumeric.ReplaceAllString(receipt.Retailer, "")))

	// Add points if total has no cents
	if receipt.Total%100 == 0 {
		points += 50
	}

	// Add points if total is a multiple of 0.25
	if receipt.Total%25 == 0 {
		points += 25
	}

	// Add points for every two items
	points += int64(len(receipt.Items) / 2 * 5)

	// Add points for eligible items
	for _, item := range receipt.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			points += int64(math.Ceil(float64(item.Price) / 500))
		}
	}

	// Add points for odd purchase date
	if receipt.PurchaseDate.Day()%2 == 1 {
		points += 6
	}

	// Add points for time of purchase
	if receipt.PurchaseTime.Before(fourPM) && receipt.PurchaseTime.After(twoPM) {
		points += 10
	}

	return Points{
		ID:     "",
		Points: points,
	}
}
