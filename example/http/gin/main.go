package main

import (
	"fmt"
	"log"

	"github.com/anthony-dong/go-sdk/commons"

	"github.com/anthony-dong/go-sdk/commons/codec/http_codec"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/api/v1/:compress", func(context *gin.Context) {
		resp := map[string]interface{}{
			"data": commons.NewString('a', 1024*8),
		}
		data := commons.ToJsonString(resp)
		context.Writer.Header().Set("Transfer-Encoding", "chunked")
		compress := context.Param("compress")
		if !http_codec.CheckAcceptEncoding(context.Request.Header, compress) {
			compress = ""
		}
		if err := http_codec.EncodeHttpBody(context.Writer, context.Writer.Header(), []byte(data), compress); err != nil {
			context.JSON(500, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
	})
	router.POST("/api/v1/test/:user", func(context *gin.Context) {

	})
	fmt.Println("Listen: http://localhost:8080\nTest: http://localhost:8080/api/v1/gzip")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
