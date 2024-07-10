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
	router := createHttpServer()
	err := router.Run("localhost:8080")
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

	router.GET("/get/:x", createShortLink)
	router.GET("/:x", redirect)
	return router
}

func redirect(c *gin.Context) {
	x := c.Params[0].Value
	url, err := D.GetLongUrl(x, db.Storage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}
	c.Redirect(http.StatusFound, url)
}

func createShortLink(c *gin.Context) {
	x := c.Params[0].Value
	url, err := D.GetShortUrl(x, db.Storage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}
