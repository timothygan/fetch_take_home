package receipt

import "time"

// Receipt
// Retailer: The name of the retailer or store the receipt is from.
// PurchaseDate: The date of the purchase printed on the receipt (YYYY-MM-DD).
// PurchaseTime: The time of the purchase printed on the receipt (24-hour format).
// Items: List of items purchased.
// Total: The total amount paid on the receipt.
type Receipt struct {
	Retailer     string    `json:"retailer" binding:"required"`
	PurchaseDate time.Time `json:"purchaseDate" binding:"required"`
	PurchaseTime time.Time `json:"purchaseTime" binding:"required"`
	Items        []Item    `json:"items" binding:"required"`
	Total        string    `json:"total" binding:"required"`
}

// ReceiptDTO - Data Transfer Object for a receipt
type ReceiptDTO struct {
	Retailer     string `json:"retailer" binding:"required"`
	PurchaseDate string `json:"purchaseDate" binding:"required"`
	PurchaseTime string `json:"purchaseTime" binding:"required"`
	Items        []Item `json:"items" binding:"required"`
	Total        string `json:"total" binding:"required"`
}

// Item
// ShortDescription: The Short Product Description for the item.
// Price: The total price paid for this item.
type Item struct {
	ShortDescription string `json:"shortDescription" binding:"required"`
	Price            string `json:"price" binding:"required"`
}
