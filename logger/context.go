package logger

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

var ctxKeyContextID = &struct{ name string }{"context_id"}

// UnwrapContext is to get the actual log context from a custom context (gin-kit context)
// Override this function at the startup to avoid race condition
var UnwrapContext = func(ctx context.Context) context.Context {
	return ctx
}

// GetID returns the context id from context
func GetID(ctx context.Context) string {
	value, _ := UnwrapContext(ctx).Value(ctxKeyContextID).(string)
	return value
}

// SaveID saves the context id to context variable and custom context store
// Return a copy of parent in which the value associated with key is
func SaveID(ctx context.Context, contextID string) context.Context {
	ctx = context.WithValue(UnwrapContext(ctx), ctxKeyContextID, contextID)
	ctx = WithContextualValues(ctx, "request_id", contextID)
	return ctx
}

// NewID creates a new id for logging
func NewID() string {
	return strings.Replace(uuid.NewString(), "-", "", -1)
}
