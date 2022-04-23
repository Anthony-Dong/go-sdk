package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	if err := http.ListenAndServe(":8888", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		for key, value := range request.Header {
			fmt.Printf("header: %v, value: %#v\n", key, value)
		}
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		fmt.Printf("body: %v\n", string(body))

		fmt.Printf("url: %v\n", request.URL.String())
		fmt.Printf("method: %v\n", request.Method)
	})); err != nil {
		panic(err)
	}
}
