package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/anthony-dong/go-sdk/commons"

	"github.com/anthony-dong/go-sdk/commons/codec/http_codec"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(func(context *gin.Context) {
		header := context.Writer.Header()
		header.Set("Access-Control-Allow-Credentials", "true")
		header.Set("Access-Control-Allow-Methods", "*")
		header.Set("Access-Control-Allow-Headers", "*")
		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Content-Security-Policy", "upgrade-insecure-requests")
	})
	router.OPTIONS("/*pattern", func(context *gin.Context) {

	})
	router.GET("/api/chunked", func(context *gin.Context) {
		getInt := func(key string, d int) int {
			query, exist := context.GetQuery(key)
			if !exist {
				return d
			}
			parseInt, err := strconv.ParseInt(query, 10, 64)
			if err != nil {
				return d
			}
			return int(parseInt)
		}
		toJson := func(v interface{}) string {
			indent, _ := json.MarshalIndent(v, "", "\t")
			return string(indent)
		}
		s1 := getInt("s1", 5)
		s2 := getInt("s2", 1)
		sleep := getInt("sleep", 100)

		context.Writer.Header().Set("Content-Type", "application/json")
		for x := 0; x < s1; x++ {
			if _, err := context.Writer.Write([]byte(toJson(map[string]interface{}{
				"value": s2,
			}))); err != nil {
				context.JSON(400, map[string]interface{}{
					"error": err,
				})
			}
			context.Writer.Flush()
			time.Sleep(time.Duration(sleep) * time.Millisecond)
		}
	})
	router.GET("/api/v1/:compress", func(context *gin.Context) {
		resp := map[string]interface{}{
			"data": commons.NewString('a', 1024),
		}
		data := commons.ToJsonString(resp)
		context.Writer.Header().Set("Content-Type", "application/json")
		context.Writer.Header().Set("Transfer-Encoding", "chunked")
		compress := context.Param("compress")
		if !http_codec.CheckAcceptEncoding(context.Request.Header, compress) {
			compress = ""
		}
		http.SetCookie(context.Writer, &http.Cookie{
			Name:  "c1",
			Value: "v1",
		})
		http.SetCookie(context.Writer, &http.Cookie{
			Name:  "c2",
			Value: "v2",
		})
		if err := http_codec.EncodeHttpBody(context.Writer, context.Writer.Header(), []byte(data), compress); err != nil {
			context.JSON(400, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		context.Writer.Flush()
	})
	router.POST("/endpoint/gencode", debug())
	router.POST("/byteapi/codegen/java/thrift/rpc_maven", debug())
	fmt.Println("sudo tcpdump -i lo0 'port 8080' -X -n | bin/gtool tcpdump")
	fmt.Println("Listen: http://localhost:8080\nTest: http://localhost:8080/api/v1/gzip")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func debug() gin.HandlerFunc {
	return func(context *gin.Context) {
		all, err := ioutil.ReadAll(context.Request.Body)
		if err != nil {
			panic(err)
		}
		if err := ioutil.WriteFile("/Users/bytedance/data/test_curl.json", all, 0644); err != nil {
			panic(err)
		}
		decoder := json.NewDecoder(bytes.NewBuffer(all))
		data := make(map[string]interface{}, 0)
		if err := decoder.Decode(&data); err != nil {
			context.JSON(500, map[string]interface{}{
				"err": err.Error(),
			})
			return
		}
		fmt.Println(len(data))
		fmt.Println(commons.ToJsonString(data))
		context.JSON(200, data)
	}
}
