package main

import (
	"fmt"

	"github.com/robertkrimen/otto"
)

func main() {

	userScript := `
var data = Mock.mock({
        'list|1-10': [{
            'id|+1': 1
        }]
    });

if (params.q1=='10'){
	mockJson={
		'data':{
			'k1':'v10'
		}	
	}
}

if (params.q1=='1'){
	mockJson={
		'data':{
			'k1': data
		}	
	}
}
	`
	vm := otto.New()
	vm.Set("params", map[string]string{
		"q1": "1",
	})
	vm.Set("header", map[string]string{
		"h1": "1",
	})
	vm.Set("mockJson", map[string]string{})
	vm.Set("mockJson_json", "") // to json string, 最后获取
	if _, err := vm.Run(appendJs(userScript)); err != nil {
		panic(err)
	} else {
		data, _ := vm.Get("mockJson_json")
		fmt.Println(data.String())
	}

}

func appendJs(js string) string {
	return `` + `
` + js + `
mockJson_json=JSON.stringify(mockJson)
`
}

//func toJson(data otto.Value) string {
//	keys := data.Object().Keys()
//	for _, elem := range keys {
//		get, err := data.Object().Get(elem)
//	}
//	return keys
//}
