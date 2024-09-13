package receipts

import (
	"errors"
)

var (
	ErrReceiptNotFound = errors.New("receipt not found")
	ErrReceiptCreate   = errors.New("failed to create receipt")
)
