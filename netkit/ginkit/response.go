package ginkit

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"hungdv136/gokit/logger"
	"hungdv136/gokit/netkit"
)

// SendJSON sends JSON
func SendJSON(ctx *gin.Context, statusCode int, verdict string, message string, data interface{}) {
	ctx.Set("status_code", statusCode)
	ctx.Set("verdict", verdict)
	ctx.Set("message", message)
	ctx.Set("data", data)

	ctx.JSON(statusCode, gin.H{
		"verdict": verdict,
		"message": message,
		"data":    data,
		"time":    time.Now().Format(time.RFC3339),
	})
}

// SendError sends error
func SendError(ctx *gin.Context, _ error) {
	message := "Unexpected error. Error ID: " + logger.GetID(ctx)
	SendJSON(ctx, http.StatusInternalServerError, netkit.VerdictFailure, message, struct{}{})
}

// SendSuccess sends success response
func SendSuccess(ctx *gin.Context, message string, data interface{}) {
	SendJSON(ctx, http.StatusOK, netkit.VerdictSuccess, message, data)
}

// SendInvalidParameters sends invalid parameters
func SendInvalidParameters(ctx *gin.Context, names []string) {
	data := map[string]interface{}{"invalid_parameters": names}
	SendJSON(ctx, http.StatusBadRequest, netkit.VerdictInvalidParameters, "invalid parameters", data)
}

// AbortJSON abort with JSON
func AbortJSON(ctx *gin.Context, statusCode int, verdict string, message string, data interface{}) {
	ctx.Abort()
	SendJSON(ctx, statusCode, verdict, message, data)
}
