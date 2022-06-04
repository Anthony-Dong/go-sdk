package main

import (
	"fmt"
	"io/ioutil"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func main() {
	file, err := ioutil.ReadFile("/Users/bytedance/data/proto/test.json")
	if err != nil {
		panic(err)
	}
	set := descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(file, &set); err != nil {
		panic(err)
	}
	marshal, err := protojson.Marshal(&set)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(marshal))
}
