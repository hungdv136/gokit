package auth

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestClaims(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	require.Nil(t, GetUserClaims(ctx))

	expectedClaims := &UserClaims{
		UserID: uuid.NewString(),
	}
	ctx = SaveUserClaims(ctx, expectedClaims)
	require.Equal(t, expectedClaims, GetUserClaims(ctx))
}
