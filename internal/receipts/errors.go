package receipts

import (
	"errors"
)

var (
	ErrReceiptNotFound = errors.New("No receipt found for that id")
	ErrReceiptInvalid  = errors.New("The receipt is invalid")
)
