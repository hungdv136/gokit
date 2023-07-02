package netkit

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestParseResponse(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("struct", func(t *testing.T) {
		t.Parallel()

		type TestStruct struct {
			Field string `json:"field"`
		}

		value := uuid.NewString()
		res, err := ParseResponse[TestStruct](ctx, &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(fmt.Sprintf(`{"field": "%s"}`, value))),
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, 200, res.StatusCode)
		require.Equal(t, value, res.Body.Field)
	})

	t.Run("map", func(t *testing.T) {
		t.Parallel()

		value := uuid.NewString()
		res, err := ParseResponse[map[string]interface{}](ctx, &http.Response{
			StatusCode: 400,
			Body:       io.NopCloser(strings.NewReader(fmt.Sprintf(`{"field": "%s"}`, value))),
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, 400, res.StatusCode)
		require.Equal(t, value, res.Body["field"])
	})
}
