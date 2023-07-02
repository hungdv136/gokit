package testkit

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hungdv136/gokit/netkit"
	"github.com/hungdv136/gokit/types"
	"github.com/hungdv136/gokit/util"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// TestCase is a test case for HTTP request
type TestCase struct {
	Name      string
	Request   *http.Request
	Assertion *Assertion
}

type Assertion struct {
	StatusCode int
	Verdict    string
}

// ResponseRecorder wraps recorder to support CloseNotify
type ResponseRecorder struct {
	*httptest.ResponseRecorder
	closeChannel chan bool
}

// CloseNotify waits for closed message
// This is required method for go-gin framework
func (r *ResponseRecorder) CloseNotify() <-chan bool {
	return r.closeChannel
}

// NewResponseRecorder returns a new instance
func NewResponseRecorder() *ResponseRecorder {
	return &ResponseRecorder{httptest.NewRecorder(), make(chan bool, 1)}
}

// NewTestCase returns a new HTTP test case. Panic if error
func NewTestCase(name string, method string, path string, expectedStatus int, expectedVerdict string) *TestCase {
	r, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}

	tc := &TestCase{
		Name:    name,
		Request: r,
		Assertion: &Assertion{
			StatusCode: expectedStatus,
			Verdict:    expectedVerdict,
		},
	}

	return tc
}

// NewUploadTestCase creates new upload test case
func NewUploadTestCase(name string, path string, file []byte, fields map[string]string, expectedStatus int, expectedVerdict string) *TestCase {
	ctx := context.Background()
	r, err := netkit.NewUploadRequest(ctx, path, file, fields)
	if err != nil {
		panic(err)
	}

	tc := &TestCase{
		Name:    name,
		Request: r,
		Assertion: &Assertion{
			StatusCode: expectedStatus,
			Verdict:    expectedVerdict,
		},
	}

	return tc
}

// WithQuery adds queries parameter. Support map and struct
func (tc *TestCase) WithQuery(queries interface{}) *TestCase {
	var m types.Map
	switch v := queries.(type) {
	case types.Map:
		m = v
	case map[string]interface{}:
		m = types.Map(v)
	default:
		var err error
		m, err = types.CreateMapFromStruct(queries)
		if err != nil {
			panic("unsupported data " + err.Error())
		}
	}

	q := tc.Request.URL.Query()
	for key, val := range m {
		q.Add(key, util.ToString(val))
	}
	tc.Request.URL.RawQuery = q.Encode()

	return tc
}

// WithBody sets request body
func (tc *TestCase) WithBody(body interface{}) *TestCase {
	data, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	r := bytes.NewReader(data)
	rc := io.NopCloser(r)
	snapshot := *r

	tc.Request.Body = rc
	tc.Request.ContentLength = int64(r.Len())
	tc.Request.GetBody = func() (io.ReadCloser, error) {
		r := snapshot
		return io.NopCloser(&r), nil
	}

	return tc
}

// WithHeader adds header
func (tc *TestCase) WithHeader(k, v string) *TestCase {
	tc.Request.Header.Add(k, v)
	return tc
}

// WithToken sets bearer token
func (tc *TestCase) WithToken(token string) *TestCase {
	return tc.WithHeader(netkit.HeaderAuthorization, "Bearer "+token)
}

// TestGin executes test case with gin engine
func TestGin[Body any](t testing.TB, tc *TestCase, engine *gin.Engine) *netkit.Response[netkit.InternalBody[Body]] {
	ctx := context.Background()
	recorder := NewResponseRecorder()
	engine.ServeHTTP(recorder, tc.Request)
	result := recorder.Result()
	defer result.Body.Close()

	require.Equal(t, tc.Assertion.StatusCode, result.StatusCode)

	res, err := netkit.ParseResponse[netkit.InternalBody[Body]](ctx, result)
	require.NoError(t, err)
	require.Equal(t, tc.Assertion.Verdict, res.Body.Verdict)

	return res
}
