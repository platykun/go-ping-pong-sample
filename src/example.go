package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/ping/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.String(http.StatusOK, "Hello %s", id)
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
