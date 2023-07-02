package ginkit

import (
	"context"
	"net/http"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/hungdv136/gokit/logger"

	"github.com/gin-gonic/gin"
)

var setupTimes = int32(0)

// OptionFunc defines server option
type OptionFunc func(s *http.Server)

// Start starts server
func Start(ctx context.Context, engine *gin.Engine, address string, options ...OptionFunc) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	srv := &http.Server{
		Addr:              address,
		Handler:           engine,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadTimeout:       30 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()

		logger.Info(ctx, "shutting down server")
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error(ctx, err)
		}
	}()

	logger.Info(ctx, "starting server", address)
	return srv.ListenAndServe()
}

// Setup one time setup. Must be called at the startup
func Setup() {
	if v := atomic.AddInt32(&setupTimes, 1); v > 1 {
		panic("This should be called at startup")
	}

	// TODO: This is opinionated solution to register UnwrapContext
	defaultFunc := logger.UnwrapContext
	logger.UnwrapContext = func(ctx context.Context) context.Context {
		if c, ok := ctx.(*gin.Context); ok {
			return c.Request.Context()
		}

		return defaultFunc(ctx)
	}

	// This is to return invalid field name as json tag instead of Golang's struct field name
	RegisterTagNameFunc()
}
