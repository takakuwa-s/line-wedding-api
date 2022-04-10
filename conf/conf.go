package conf

import (
	"context"

	"go.uber.org/zap"
)

var (
	Log, _ = zap.NewProduction()
	Ctx = context.Background()
)