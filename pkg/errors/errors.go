package errors

import (
	"errors"
	"fmt"
)

var (
	ErrURLNotFound  = errors.New("url not found")
	ErrNotUniqueURL = errors.New("url not unique error")
	ErrDeleted      = errors.New("deleted")
)

type NotUniqueURLErr struct {
	URLID       string
	OriginalURL string
	Err         error
}

func NewNotUniqueURLErr(urlID, originalURL string, err error) error {
	return &NotUniqueURLErr{
		URLID:       urlID,
		OriginalURL: originalURL,
		Err:         err,
	}
}

func (e *NotUniqueURLErr) Error() string {
	return fmt.Sprintf("url not unique: urlID %v, originalURL: %v, error: %v ", e.URLID, e.OriginalURL, e.Err)
}
