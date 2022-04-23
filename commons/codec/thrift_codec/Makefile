.PHONY: gen

gen: ## 注意构建成功后可能需要手动修改代码！因为thrift生成的代码可能存在问题
	thrift -r --gen \
	go:thrift_import=github.com/apache/thrift/lib/go/thrift,package_prefix=github.com/anthony-dong/go-tool/command/codec/thrift_codec/test/ \
	-out ./test \
	"${GOPATH}/src/github.com/anthony-dong/go-tool/command/codec/thrift_codec/test/ThriftStruct.thrift"
	rm -rf test/thriftstruct/test_service-remote