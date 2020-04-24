package testutil

import (
	"context"
	"github.com/NYTimes/gizmo/server/kit"
	"github.com/go-kit/kit/log"
)

func NewTestContext() context.Context {
	logger := log.NewNopLogger()
	return kit.SetLogger(context.Background(), logger)
}