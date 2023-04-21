package logs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStdOutLogger(t *testing.T) {
	//config := zap.NewProductionConfig()
	//config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	//config.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	//config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	//config.Encoding = "console"
	//
	//logger, _ := config.Build()
	//
	////zap.New(core)
	////logger, _ := zap.NewProduction()
	////defer logger.Sync() // flushes buffer, if any
	////
	//sugar := logger.Sugar()
	//
	//sugar.Infof("Failed to fetch URL: %s", "xxxx")
}

func TestName(t *testing.T) {
	Errorf("err")
	Infof("info")
	SetLevel(LevelInfo)
	Debugf("11111")
	t.Run("", func(t *testing.T) {
		SetLevel(LevelDebug)
		assert.Equal(t, IsLevel(LevelDebug), true)
		assert.Equal(t, IsLevel(LevelInfo), true)
	})
	t.Run("", func(t *testing.T) {
		SetLevel(LevelFatal)
		assert.Equal(t, IsLevel(LevelFatal), true)
		assert.Equal(t, IsLevel(LevelInfo), false)
	})

	SetLevel(LevelDebug)
	CtxDebugf(SetLogId(context.Background(), "1111"), "data: %s", "1")
	CtxDebugf(context.WithValue(context.Background(), struct {
	}{}, "xxxxx"), "data: %s", "1")

	assert.Equal(t, IsLevel(LevelDebug), true)
	assert.Equal(t, IsLevel(LevelError), true)

	SetLevel(LevelFatal)
	assert.Equal(t, IsLevel(LevelFatal), true)
	assert.Equal(t, IsLevel(LevelError), false)
}
