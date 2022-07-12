package models

type ShortenRequest struct {
	URL string `json:"url" valid:"url,required"`
}

type ShortenReply struct {
	ShortenURLResult string `json:"result"`
}

type UserURL struct {
	ShortURL    string
	OriginalURL string
}
