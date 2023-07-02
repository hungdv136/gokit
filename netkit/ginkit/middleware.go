package ginkit

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"hungdv136/gokit/logger"
	"hungdv136/gokit/netkit"
	"hungdv136/gokit/netkit/auth"
	"hungdv136/gokit/types"
)

// RequestIDMiddleware adds X-Request-ID value to request, response and save to context variable
func RequestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Request.Header.Get(netkit.HeaderRequestID)
		if len(id) == 0 {
			id = logger.NewID()
			ctx.Request.Header.Set(netkit.HeaderRequestID, id)
		}

		wrappedCtx := logger.SaveID(ctx, id)
		ctx.Request = ctx.Request.WithContext(wrappedCtx)
		ctx.Writer.Header().Set(netkit.HeaderRequestID, id)
		ctx.Next()
	}
}

// RequestTimeMiddleware logs request time
func RequestTimeMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		ctx.Next()
		runTime := time.Since(startTime)

		logger.Fields(ctx, "remote_addr", ctx.Request.RemoteAddr,
			"method", ctx.Request.Method,
			"path", ctx.Request.URL.Path,
			"route_path", ctx.FullPath(),
			"status_code", ctx.Writer.Status(),
			"response_verdict", ctx.GetString("verdict"),
			"response_message", ctx.GetString("message"),
			"duration_ns", runTime.Nanoseconds(),
		).Info(ctx)
	}
}

// TracingMiddleware adds tracing into gin's middlewares by using otelgin.Middleware function
func TracingMiddleware(service string, opts ...otelgin.Option) gin.HandlerFunc {
	return otelgin.Middleware(service, opts...)
}

// Recovery returns a middleware for a given writer that recovers from any panics and calls the provided handle func to handle it
func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if rErr := recover(); rErr != nil {
				err, ok := rErr.(error)
				if !ok {
					err = fmt.Errorf("invalid recover error %v", rErr)
				}

				logger.Fields(ctx, "error.stack", debug.Stack()).Error(ctx, "panic recovered:", err)

				// If the connection is dead, we can't write a status to it.
				if isBrokenPipeError(rErr) {
					_ = ctx.Error(err)
					ctx.Abort()
				} else {
					ctx.Abort()
					SendError(ctx, err)
				}
			}
		}()

		ctx.Next()
	}
}

// MaxBody checks maxbody payload
func MaxBody(maxBytes int64) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.ContentLength > maxBytes {
			msg := fmt.Sprintf("body size exceeds %d bytes", maxBytes)
			logger.Warn(ctx, "invalid token")
			AbortJSON(ctx, http.StatusBadRequest, netkit.VerdictInvalidParameters, msg, struct{}{})
			return
		}

		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxBytes)
		ctx.Next()
	}
}

// Authenticate authenticate user
func Authenticate(verifier auth.Verifier) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get(netkit.HeaderAuthorization)
		token = strings.TrimSpace(strings.TrimPrefix(token, netkit.TokenTypeBearer))
		if len(token) == 0 {
			logger.Warn(ctx, "missing token")
			AbortJSON(ctx, http.StatusUnauthorized, netkit.VerdictMissingAuthentication, "missing authorization header", types.Map{})
			return
		}

		claims, err := verifier.Verify(ctx, token)
		if err != nil {
			logger.Warn(ctx, "invalid token")
			AbortJSON(ctx, http.StatusUnauthorized, netkit.VerdictInvalidToken, "invalid token", types.Map{})
			return
		}

		ctx.Set(gin.AuthUserKey, claims)
		wrappedCtx := auth.SaveUserClaims(ctx.Request.Context(), claims)
		wrappedCtx = logger.WithContextualValues(wrappedCtx, "user_id", claims.UserID)
		ctx.Request = ctx.Request.WithContext(wrappedCtx)
	}
}

// Check for a broken connection, as it is not really a condition that warrants a panic stack trace
func isBrokenPipeError(err interface{}) bool {
	ne, ok := err.(*net.OpError)
	if !ok {
		return false
	}

	se := &os.SyscallError{}
	if errors.As(ne, &se) {
		if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
			return true
		}
	}

	return false
}
