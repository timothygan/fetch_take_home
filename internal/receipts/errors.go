package receipts

import (
	"errors"
)

var (
	ErrReceiptNotFound = errors.New("No receipt found for that id")
	ErrReceiptCreate   = errors.New("The receipt is invalid")
)
