.PHONY: all  protoc
all: protoc

## install  protoc & protoc-gen-go & protoc-gen-go-grpc
## wget https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-osx-x86_64.zip
## go get -v google.golang.org/protobuf/cmd/protoc-gen-go
## go get -v google.golang.org/grpc/cmd/protoc-gen-go-grpc

protoc: ## 编译上面的'test.proto'文件
	protoc \
    --experimental_allow_proto3_optional \
    --proto_path=. \
    --plugin=protoc-gen-go=${HOME}/go/bin/protoc-gen-go \
    --go_opt=Mtest.proto=github.com/anthony-dong/go-tool/internal/example/protobuf/test \
    --go_out=${HOME}/go/src \
    --plugin=protoc-gen-go-grpc=${HOME}/go/bin/protoc-gen-go-grpc \
    --go-grpc_opt=Mtest.proto=github.com/anthony-dong/go-tool/internal/example/protobuf/test \
    --go-grpc_out=${HOME}/go/src \
    test.proto