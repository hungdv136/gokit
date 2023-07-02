package netkit

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hungdv136/rio"
	"github.com/stretchr/testify/require"
)

func TestSendJSON(t *testing.T) {
	t.Parallel()

	server := rio.NewLocalServerWithReporter(t)

	ctx := context.Background()
	type (
		ResponseBody struct {
			Field string `json:"field"`
		}

		RequestBody struct {
			ID string `json:"id"`
		}
	)

	resBody := ResponseBody{
		Field: uuid.NewString(),
	}

	reqBody := RequestBody{
		ID: uuid.NewString(),
	}

	require.NoError(t, rio.NewStub().
		For("POST", rio.Contains("animal")).
		WithRequestBody(rio.BodyJSONPath("$.id", rio.EqualTo(reqBody.ID))).
		WillReturn(rio.JSONResponse(resBody)).
		Send(ctx, server))

	res, err := SendJSON[ResponseBody](ctx, "POST", server.GetURL(ctx)+"/animal", reqBody)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, 200, res.StatusCode)
	require.Equal(t, resBody.Field, res.Body.Field)
}

func TestGetRequest(t *testing.T) {
	t.Parallel()

	server := rio.NewLocalServerWithReporter(t)

	ctx := context.Background()
	type ResponseBody struct {
		Field string `json:"field"`
	}

	resBody := ResponseBody{
		Field: uuid.NewString(),
	}

	require.NoError(t, rio.NewStub().
		For("GET", rio.Contains("animal")).
		WillReturn(rio.JSONResponse(resBody)).
		Send(ctx, server))

	res, err := Get[ResponseBody](ctx, server.GetURL(ctx)+"/animal")
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, 200, res.StatusCode)
	require.Equal(t, resBody.Field, res.Body.Field)
}
