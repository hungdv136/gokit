package testkit

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hungdv136/gokit/types"
	"github.com/hungdv136/rio"
	"github.com/stretchr/testify/require"
)

func TestWithQuery(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	server := rio.NewLocalServerWithReporter(t)

	req := struct {
		F1 string `json:"f1"`
		F2 string `json:"f2"`
	}{
		F1: uuid.NewString(),
		F2: uuid.NewString(),
	}

	require.NoError(t, rio.NewStub().
		For("GET", rio.Contains("/animal")).
		WithQuery("f1", rio.Contains(req.F1)).
		WithQuery("f2", rio.Contains(req.F2)).
		WillReturn(rio.JSONResponse(types.Map{"verdict": "success"})).
		Send(ctx, server))

	tc := NewTestCase(t.Name(), "GET", server.GetURL(ctx)+"/animal", 200, "success").WithQuery(req)
	Execute[types.Map](t, tc)
}

func TestWithBody(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	server := rio.NewLocalServerWithReporter(t)

	req := struct {
		F1 string `json:"f1"`
		F2 string `json:"f2"`
	}{
		F1: uuid.NewString(),
		F2: uuid.NewString(),
	}

	require.NoError(t, rio.NewStub().
		For("POST", rio.Contains("/animal")).
		WithRequestBody(rio.BodyJSONPath("$.f1", rio.Contains(req.F1))).
		WithRequestBody(rio.BodyJSONPath("$.f2", rio.Contains(req.F2))).
		WillReturn(rio.JSONResponse(types.Map{"verdict": "success"})).
		Send(ctx, server))

	tc := NewTestCase(t.Name(), "POST", server.GetURL(ctx)+"/animal", 200, "success").WithBody(req)
	Execute[types.Map](t, tc)
}
