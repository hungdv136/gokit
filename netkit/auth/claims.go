package auth

import "context"

// key for saving user claims to the context
var ctxKeyUserClaims = &struct{ name string }{"user_claims"}

// SaveUserClaims saves authorized user claims to the context
func SaveUserClaims(ctx context.Context, claims *UserClaims) context.Context {
	return context.WithValue(ctx, ctxKeyUserClaims, claims)
}

// GetUserClaims returns the authorized user's claims from context
func GetUserClaims(ctx context.Context) *UserClaims {
	claims := ctx.Value(ctxKeyUserClaims)
	if claims == nil {
		return nil
	}

	return claims.(*UserClaims)
}
