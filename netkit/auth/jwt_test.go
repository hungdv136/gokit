package auth

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUserJWT(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	signer, err := NewRandomJwt()
	require.NoError(t, err)

	r := &UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:       uuid.NewString(),
			Issuer:   uuid.NewString(),
			Audience: []string{uuid.NewString()},
		},
		UserID: uuid.NewString(),
		Email:  uuid.NewString(),
		Name:   uuid.NewString(),
	}
	tokenString, err := signer.Sign(ctx, r)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)
	require.True(t, len(tokenString) <= 1024)

	claims, err := signer.Verify(context.Background(), tokenString)
	require.NoError(t, err)
	require.NotNil(t, claims)
	require.NotEmpty(t, claims.ID)
	require.NoError(t, err)
	require.NotNil(t, claims)
	require.Equal(t, r.UserID, claims.UserID)
	require.Equal(t, r.Audience, claims.Audience)
	require.Equal(t, r.Issuer, claims.Issuer)
}
