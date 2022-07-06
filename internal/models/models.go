package models

type ShortenRequest struct {
	URL string `json:"url" valid:"url,required"`
}

type ShortenReply struct {
	ShortenURLResult string `json:"result"`
}

type UserURL struct {
	ID          int    `db:"id"`
	UserID      string `json:"user_id" db:"user"`
	ShortURL    string `json:"short_url" db:"short_url"`
	OriginalURL string `json:"original_url" db:"original_url"`
}
