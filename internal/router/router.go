package router

import (
	"github.com/ChristinaFomenko/shortener/configs"
	"github.com/ChristinaFomenko/shortener/internal/app/storage"
	"github.com/ChristinaFomenko/shortener/internal/handlers"
	"github.com/ChristinaFomenko/shortener/internal/middlewares"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func Router(c configs.AppConfig, s storage.Repository) *gin.Engine {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.BestSpeed, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	router.Use(middlewares.SessionAuthMiddleware(c))

	handler := handlers.Handler{
		Config:  c,
		Storage: s,
	}
	router.POST("/", handler.Shorten)
	router.GET("/{id}", handler.Expand)
	router.POST("/api/shorten", handler.APIJSONShorten)
	router.GET("/api/user/urls", handler.GetList)
	router.GET("/ping", handler.Ping)
	router.POST("/api/shorten/batch", handler.BatchShortenHandler)
	//})
	return router
}
