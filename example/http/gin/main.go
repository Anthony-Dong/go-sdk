package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/api/v1/:user", func(context *gin.Context) {

	})
	router.POST("/api/v1/test/:user", func(context *gin.Context) {

	})

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
