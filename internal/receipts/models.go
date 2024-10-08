package receipts

import "time"

// Receipt
// ID: UUID of the receipt
// Retailer: The name of the retailer or store the receipt is from.
// PurchaseDate: The date of the purchase printed on the receipt (YYYY-MM-DD).
// PurchaseTime: The time of the purchase printed on the receipt (24-hour format).
// Items: List of items purchased.
// Total: The total amount paid on the receipt.
type Receipt struct {
	ID           string    `json:"id"`
	Retailer     string    `json:"retailer"`
	PurchaseDate time.Time `json:"purchaseDate"`
	PurchaseTime time.Time `json:"purchaseTime"`
	Items        []Item    `json:"items"`
	Total        int64     `json:"total"`
}

// Item
// ShortDescription: The Short Product Description for the item.
// Price: The total price paid for this item in cents
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            int64  `json:"price"`
}

// ReceiptDTO - Data Transfer Object for a receipt
type ReceiptDTO struct {
	Retailer     string    `json:"retailer" binding:"required"`
	PurchaseDate string    `json:"purchaseDate" binding:"required"`
	PurchaseTime string    `json:"purchaseTime" binding:"required"`
	Items        []ItemDTO `json:"items" binding:"required"`
	Total        string    `json:"total" binding:"required"`
}

// ItemDTO
// ShortDescription: The Short Product Description for the item.
// Price: The total price paid for this item.
type ItemDTO struct {
	ShortDescription string `json:"shortDescription" binding:"required"`
	Price            string `json:"price" binding:"required"`
}

// Points
// ID: The ID of the receipt
// Points: The number of points awarded
type Points struct {
	ID     string `json:"id"`
	Points int64  `json:"points"`
}

// CreateResponse
// id: The ID of the receipt
type CreateResponse struct {
	ID string `json:"id"`
}

// PointsResponse
// points: The number of points awarded
type PointsResponse struct {
	Points int64 `json:"points"`
}
