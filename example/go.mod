module github.com/anthony-dong/go-sdk/example

go 1.13

require (
	github.com/anthony-dong/go-sdk v0.0.0-20220423155222-042777f402c5
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/epiclabs-io/diff3 v0.0.0-20181217103619-05282cece609
	github.com/gin-gonic/gin v1.7.7
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da
	github.com/jhump/protoreflect v1.12.0
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/kr/pty v1.1.1
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/robertkrimen/otto v0.0.0-20211024170158-b87d35c0b86f
	github.com/stretchr/testify v1.7.1
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.28.0
)

replace github.com/anthony-dong/go-sdk => ../
