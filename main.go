package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	DB "shortUrl/db"
	D "shortUrl/domain"
)

var db = DB.NewService(DB.NewSQLiteRepository())

func main() {
	err := db.Storage.CreateTable()
	if err != nil {
		DB.NewSQLiteRepository()
		return
	}
	router := createHttpServer()
	err = router.Run("localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func createHttpServer() *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.POST("/create", createShortLink)
	router.GET("/:x", redirect)
	return router
}

func redirect(c *gin.Context) {
	x := c.Params[0].Value
	url, err := D.GetLongUrl(x, db.Storage)
	if url == "" {
		c.JSON(http.StatusNotFound, gin.H{"err": "url not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	c.Redirect(http.StatusFound, url)
}

type CreateShortLinkRequest struct {
	Url string `json:"url"`
}

func createShortLink(c *gin.Context) {
	var json CreateShortLinkRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	url, err := D.GetShortUrl(json.Url, db.Storage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}
