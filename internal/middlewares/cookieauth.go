package middlewares

import (
	"github.com/ChristinaFomenko/shortener/configs"
	"github.com/ChristinaFomenko/shortener/internal/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"net/url"
)

func SessionAuthMiddleware(conf configs.AppConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("user_id")

		if cookie == "" || err != nil {
			encryptedID, err := utils.Encrypt(uuid.New().String(), conf.AuthKey)
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}

			ctx.Request.AddCookie(&http.Cookie{
				Name:  "session",
				Value: url.QueryEscape(encryptedID),
			})

			ctx.SetCookie("user_id", encryptedID, 3600, "/", conf.ServerAddress, false, false)
		}

		ctx.Next()
	}
}
