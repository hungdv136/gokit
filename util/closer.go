package util

import (
	"context"

	"github.com/hungdv136/gokit/logger"
)

// CloseSilently closes and write log if an error occurs
func CloseSilently(ctx context.Context, close func() error) {
	if err := close(); err != nil {
		logger.Error(ctx, "cannot close", err)
	}
}
