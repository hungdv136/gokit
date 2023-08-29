package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hungdv136/gokit/logger"
)

// Some default values
const (
	DefaultAlgorithm = "RS256"
)

// UserClaims is a custom claims that contains more user data beside standard claims's data
type UserClaims struct {
	jwt.RegisteredClaims
	Name          string `json:"name"`
	PictureURL    string `json:"picture_url"`
	UserID        string `json:"user_id"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

// Verifier defines interface to verify authentication token
type Verifier interface {
	Verify(ctx context.Context, tokenString string) (*UserClaims, error)
}

// Signer defines an interface to sign authentication token
type Signer interface {
	Sign(ctx context.Context, claims *UserClaims) (string, error)
}

// Jwt is to sign and verify JWT
type Jwt struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// NewJWTFromPublicPem initializes jwt verifier from public string
// An error will be returned if returned object is used to sign a JWT
func NewJWTFromPublicPem(pem string) (*Jwt, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
	if err != nil {
		return nil, err
	}

	return &Jwt{PublicKey: key}, nil
}

// NewJWTFromPrivatePem creates signer and verifier from PEM string
func NewJWTFromPrivatePem(pem string) (*Jwt, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pem))
	if err != nil {
		return nil, err
	}

	return &Jwt{PrivateKey: key, PublicKey: &key.PublicKey}, nil
}

// NewRandomJwt generates random RSA keys for RSA algorithm
// This is to simplify setup jwt verify for unit test
func NewRandomJwt() (*Jwt, error) {
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return &Jwt{PrivateKey: k, PublicKey: &k.PublicKey}, nil
}

// Sign creates a new JWT token signed by RSA method
func (a *Jwt) Sign(ctx context.Context, claims *UserClaims) (string, error) {
	if a.PrivateKey == nil {
		err := errors.New("missing private key")
		logger.Error(ctx, err)
		return "", err
	}

	token := jwt.New(jwt.GetSigningMethod(DefaultAlgorithm))
	token.Claims = claims
	signed, err := token.SignedString(a.PrivateKey)
	if err != nil {
		logger.Error(ctx, err)
		return "", err
	}

	return signed, nil
}

// Verify checks if provided token string is valid or not
// claims is always returned to respect the jwtgo's behavior
// caller should have decision depend on the error type
func (a *Jwt) Verify(ctx context.Context, tokenString string) (*UserClaims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return a.PublicKey, nil
	}

	var claims UserClaims
	if _, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc); err != nil {
		if IsExpiredJWTError(err) {
			logger.Warn(ctx, err)
			return &claims, err
		}

		logger.Error(ctx, err, tokenString)
		return &claims, err
	}

	return &claims, nil
}

// IsExpiredJWTError checks if err is JWT ValidationErrorExpired
func IsExpiredJWTError(err error) bool {
	return errors.Is(err, jwt.ErrTokenExpired)
}
