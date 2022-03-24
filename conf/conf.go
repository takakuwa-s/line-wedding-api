package conf

import (
	"go.uber.org/zap"
)

var (
	Log, _ = zap.NewProduction()
)