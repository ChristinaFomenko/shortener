package handlers

type GetUrlsReply struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
