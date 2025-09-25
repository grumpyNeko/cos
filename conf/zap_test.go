package conf

import (
	"context"
	"go.uber.org/zap"
	"testing"
)

func Test_zap(t *testing.T) {
	ctx := context.WithValue(context.TODO(), "traceid", 2)
	do(ctx)
}

func do(ctx context.Context) {
	log := Stdout.With(zap.Any("traceid", ctx.Value("traceid")))
	log.Error("Todo", zap.Any("userid", 2))
}
