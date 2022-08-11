package tcpdump

import (
	"context"
	"testing"
)

func TestContext_HandlerPacket(t *testing.T) {
	ctx := NewCtx(context.Background(), NewDefaultConfig())
	ctx.HandlerPacket(Packet{
		Src:  "127.0.0.1:8888",
		Dst:  "127.0.0.1:8889",
		Data: []byte("hello"),
	})
}
