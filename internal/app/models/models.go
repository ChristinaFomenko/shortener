package models

type OriginalURL struct {
	CorrelationID string
	URL           string
}

type UserURL struct {
	CorrelationID string
	ShortURL      string
	OriginalURL   string
	IsDeleted     bool `db:"is_deleted"`
}
