package curl

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	curl1 = `curl  --location --request POST 'http://www.baidu.com/api/v1/get_api?use_hijack=true&tag=test&device_id=69262354575' \
--header 'get: 1' \
--header 'Content-Type: application/json;charset=UTF-8' \
--data '{
    "Name": "fanhaodong.test.item",
    "MethodName": "GetItem",
    "RequestBody": "{\"id\":1,\"name\":\"hello\"}",
    "Instance": ""
}'`
	curl2 = `curl --location 'http://localhost:8888/api/v1/test/mock/rpc?device_id=6926235457512382119122&tag=&use_hijack=true'  --header 'Content-Type: application/json;charset=UTF-8' --header 'Content-Type: application/json'  --data-raw "{
    \"Name\": \"\",
    \"MethodName\": \"\",
    \"RequestBody\": \"\",
    \"Instance\": \"\",
    \"ReqTimeout\": 0,
    \"BaseData\": {
        \"LogID\": \"\",
        \"Caller\": \"\",
        \"Addr\": \"\",
        \"Client\": \"\",
        \"TrafficEnv\": {
            \"Open\": false,
            \"Env\": \"\"
        },
        \"Extra\": {
            \"\": \"\"
        }
    }
}"`

	curl3 = `curl --location --request POST 'http://localhost:8888/api/v1/test/mock/rpc?device_id=6926235457512382119122&tag=&use_hijack=true'  --header 'Content-Type: application/json'  --data-raw '{
    \"Name\": \"\",
    \"MethodName\": \"\",
    \"RequestBody\": \"\",
    \"Instance\": \"\",
    \"ReqTimeout\": 0,
    \"BaseData\": {
        \"LogID\": \"\",
        \"Caller\": \"\",
        \"Addr\": \"\",
        \"Client\": \"\",
        \"TrafficEnv\": {
            \"Open\": false,
            \"Env\": \"\"
        },
        \"Extra\": {
            \"\": \"\"
        }
    }
}'`
)

func TestParserCmd2Slice(t *testing.T) {
	t.Run("t1", func(t *testing.T) {
		for _, elem := range ParserCmd2Slice(curl1) {
			t.Log(elem)
		}
	})

	t.Run("t2", func(t *testing.T) {
		for _, elem := range ParserCmd2Slice(curl2) {
			t.Log(elem)
		}
	})
	t.Run("t3", func(t *testing.T) {
		for _, elem := range ParserCmd2Slice(curl3) {
			t.Log(elem)
		}
	})
}

func TestToHttpInfo(t *testing.T) {
	testCase := []string{curl1, curl2, curl3}
	for _, elem := range testCase {
		t.Log("==========================")
		info, err := ToHttpInfo(elem)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("url: %v\n", info.Url)
		t.Logf("method: %v\n", info.Method)
		t.Logf("header: %v\n", info.Header)
		t.Logf("body: %#v", info.Body)
	}
}

func TestName(t *testing.T) {
	assert.Equal(t, isUrl("http://127.0.0.1:8888/api/v1/test/mock/rpc?device_id=6926235457512382119122&tag=&use_hijack=true"), true)
	assert.Equal(t, isUrl("http://localhost:8888/api/v1/test/mock/rpc?device_id=6926235457512382119122&tag=&use_hijack=true"), true)
	assert.Equal(t, isUrl("localhost:8888/api/v1/test/mock/rpc?device_id=6926235457512382119122&tag=&use_hijack=true"), false)
}

func TestHttpHeader(t *testing.T) {
	header := http.Header{}
	header.Add("h1", "v1")
	header.Add("h1", "v2")
	t.Log(header)
}
