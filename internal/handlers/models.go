package handlers

type ShortenRequest struct {
	URL string `json:"url" valid:"url,required"`
}

type ShortenReply struct {
	ShortenURLResult string `json:"result"`
}

type GetUrlsReply struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type ShortenBatchRequest struct {
	CorrelationID string `json:"correlation_id" valid:"required"`
	OriginalURL   string `json:"original_url" valid:"url,required"`
}

type ShortenBatchReply struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
