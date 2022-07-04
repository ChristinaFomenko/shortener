package handlers

import "github.com/ChristinaFomenko/shortener/internal/app/models"

func toGetUrlsReply(model []models.UserURL) []GetUrlsReply {
	reply := make([]GetUrlsReply, len(model))

	for idx, m := range model {
		reply[idx] = GetUrlsReply{
			ShortURL:    m.ShortURL,
			OriginalURL: m.OriginalURL,
		}
	}

	return reply
}
