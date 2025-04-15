package common

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type StrIns interface {
	Str() string
}

func ZapU(p StrIns) zapcore.Field {
	return zap.String("U", p.Str())
}
