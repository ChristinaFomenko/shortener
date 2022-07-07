package storage

type UserURL struct {
	ID          string
	OriginalURL string
	UserID      string
}

type URLDuplicateError struct {
	URL string
}
