package netkit

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"hungdv136/gokit/logger"
)

// Response is generic response structure
// Body is data type of the response body
type Response[Body any] struct {
	StatusCode int  `json:"status_code"`
	Body       Body `json:"body"`
}

// ParseResponse parses an http response to struct
func ParseResponse[R any](ctx context.Context, r *http.Response) (*Response[R], error) {
	resp := Response[R]{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&resp.Body); err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	resp.StatusCode = r.StatusCode
	return &resp, nil
}

// InternalBody defines struct for internal response body
// This is for communicate with internal services
type InternalBody[R any] struct {
	Verdict string    `json:"verdict"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
	Data    R         `json:"data"`
}
