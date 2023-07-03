package netkit

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hungdv136/gokit/logger"
	"github.com/hungdv136/gokit/util"
)

var (
	headerIPAddress = []string{"X-Forwarded-For", "X-Real-Ip"}
	defaultClient   = &http.Client{
		Timeout:   30 * time.Second,
		Transport: NewTransport(WithIdleConnsPerHost(10)),
	}
)

// UploadOptions upload options
type UploadOptions struct {
	Method    string            // Default: POST
	FieldName string            // Default: file
	FileName  string            // Default: UUID
	Fields    map[string]string // Default: Empty
}

// SendRequest sends general request to a URL and returns HTTP response
func SendRequest(r *http.Request) (*http.Response, error) {
	if id := logger.GetID(r.Context()); len(id) > 0 {
		r.Header.Add(HeaderRequestID, id)
	}

	response, err := defaultClient.Do(r)
	if err != nil {
		logger.Error(r.Context(), err)
		return nil, err
	}

	return response, nil
}

// NewUploadRequest create a new http upload file request
func NewUploadRequest(ctx context.Context, url string, reader io.Reader, modifiers ...func(*UploadOptions)) (*http.Request, error) {
	options := &UploadOptions{
		Method:    http.MethodPost,
		FieldName: "file",
		FileName:  uuid.NewString(),
	}

	for _, modifier := range modifiers {
		modifier(options)
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(options.FieldName, options.FileName)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	if _, err := io.Copy(part, reader); err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	for key, val := range options.Fields {
		if err := writer.WriteField(key, val); err != nil {
			logger.Error(ctx, err)
			return nil, err
		}
	}

	if err := writer.Close(); err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, options.Method, url, body)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req.WithContext(ctx), nil
}

// NewJSONRequest creates new request with JSON body
func NewJSONRequest(ctx context.Context, method, url string, body interface{}) (*http.Request, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		logger.Error(ctx, "cannot marshal body", err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyBytes))
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// NewQueryRequest creates request with query strings
func NewQueryRequest(ctx context.Context, method, url string, queries map[string]interface{}) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		logger.Error(ctx, "cannot create request", err)
		return nil, err
	}

	if len(queries) > 0 {
		q := req.URL.Query()
		for key, val := range queries {
			q.Add(key, util.ToString(val))
		}
		req.URL.RawQuery = q.Encode()
	}

	return req, nil
}

// SendJSON executes request JSON as body and parses JSON response body
// Body is structure of response body
func SendJSON[Body any](ctx context.Context, method, url string, body interface{}) (*Response[Body], error) {
	req, err := NewJSONRequest(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	res, err := SendRequest(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	return ParseResponse[Body](ctx, res)
}

// Get executes request with GET method and parses JSON response body
// Body is structure of response body
func Get[Body any](ctx context.Context, url string) (*Response[Body], error) {
	req, err := NewQueryRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := SendRequest(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	return ParseResponse[Body](ctx, res)
}

// GetRemoteAddr gets IP of caller
func GetRemoteAddr(r *http.Request) string {
	for _, h := range headerIPAddress {
		ips := r.Header.Get(h)
		if len(ips) == 0 {
			continue
		}

		i := strings.LastIndex(strings.TrimSuffix(ips, ","), ",")
		return ips[i+1:]
	}

	return r.RemoteAddr
}
