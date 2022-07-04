package handlers

import "github.com/ChristinaFomenko/shortener/internal/app/models"

func toGetUrlsReply(model []models.UserURL) []models.GetUrlsReply {
	reply := make([]models.GetUrlsReply, len(model))

	for idx, m := range model {
		reply[idx] = models.GetUrlsReply{
			ShortURL:    m.ShortURL,
			OriginalURL: m.OriginalURL,
		}
	}

	return reply
}
