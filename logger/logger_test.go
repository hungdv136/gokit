package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"hungdv136/gokit/types"
)

func TestLogger(t *testing.T) {
	writer := bytes.Buffer{}
	Setup("test-logger", &writer)

	id := uuid.NewString()
	fieldValue := uuid.NewString()
	message := uuid.NewString()

	ctx := SaveID(context.Background(), id)
	Fields(ctx, "key", fieldValue).Info(ctx, message)

	output := types.Map{}
	require.NoError(t, json.Unmarshal(writer.Bytes(), &output))
	require.Equal(t, id, output["request_id"])
	require.Equal(t, fieldValue, output["key"])
	require.Equal(t, message, output["message"])
}
