package ginkit

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"hungdv136/gokit/logger"
	"hungdv136/gokit/netkit"
	"hungdv136/gokit/netkit/auth"
	"hungdv136/gokit/netkit/testkit"
	"hungdv136/gokit/types"
)

func TestRecoveryMiddleware(t *testing.T) {
	t.Parallel()

	handler := func(ctx *gin.Context) {
		p := ctx.Query("what")
		if p == "panic" {
			panic(errors.New("something wrong"))
		}

		SendSuccess(ctx, "success", types.Map{"what": p})
	}

	engine := gin.New()
	engine.Use(Recovery())
	engine.GET("/test/recovery", handler)

	t.Run("panic", func(t *testing.T) {
		t.Parallel()

		param := types.Map{"what": "panic"}
		tc := testkit.NewTestCase("panic", "GET", "/test/recovery", 500, netkit.VerdictFailure).WithQuery(param)
		res := testkit.TestGin[types.Map](t, tc, engine)
		require.Empty(t, res.Body.Data.ForceMap("data"))
	})

	t.Run("no_panic", func(t *testing.T) {
		t.Parallel()

		param := types.Map{"what": uuid.NewString()}
		tc := testkit.NewTestCase("panic", "GET", "/test/recovery", 200, netkit.VerdictSuccess).WithQuery(param)
		res := testkit.TestGin[types.Map](t, tc, engine)
		require.Empty(t, res.Body.Data.ForceMap("data"))
	})
}

func TestRequestIDMiddleware(t *testing.T) {
	t.Parallel()

	handler := func(ctx *gin.Context) {
		requestID := logger.GetID(ctx)
		headerID := ctx.Request.Header.Get(netkit.HeaderRequestID)

		if len(requestID) == 0 {
			logger.Info(ctx, "request id is empty")
			SendError(ctx, errors.New("request id is empty"))
			return
		}

		if len(headerID) > 0 && requestID != headerID {
			logger.Info(ctx, "request id does not match")
			SendError(ctx, errors.New("request id does not match"))
			return
		}

		SendSuccess(ctx, "success", types.Map{"request_id": requestID})
	}

	engine := gin.New()
	engine.Use(RequestIDMiddleware())
	engine.GET("/test", handler)

	expectedID := uuid.NewString()
	testCases := []*testkit.TestCase{
		testkit.NewTestCase("with_header", "GET", "/test", 200, netkit.VerdictSuccess).WithHeader(netkit.HeaderRequestID, expectedID),
		testkit.NewTestCase("without_header", "GET", "/test", 200, netkit.VerdictSuccess),
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			res := testkit.TestGin[types.Map](t, tc, engine)
			requestID := res.Body.Data.ForceString("request_id")
			require.NotEmpty(t, requestID)

			if tc.Name == "with_header" {
				require.Equal(t, expectedID, requestID)
			}
		})
	}
}

func TestAuthenticate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	handler := func(ctx *gin.Context) {
		claims := auth.GetUserClaims(ctx.Request.Context())
		if claims == nil {
			SendJSON(ctx, 401, netkit.VerdictInvalidToken, "missing", types.Map{})
			return
		}

		SendSuccess(ctx, "success", types.Map{"claims": claims})
	}

	jwt, err := auth.NewRandomJwt()
	require.NoError(t, err)

	engine := gin.New()
	engine.Use(Authenticate(jwt))
	engine.GET("/test", handler)

	claims := &auth.UserClaims{UserID: uuid.NewString()}
	token, err := jwt.Sign(ctx, claims)
	require.NoError(t, err)

	testCases := []*testkit.TestCase{
		testkit.NewTestCase("success", "GET", "/test", 200, netkit.VerdictSuccess).WithToken(token),
		testkit.NewTestCase("missing", "GET", "/test", 401, netkit.VerdictMissingAuthentication),
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			res := testkit.TestGin[types.Map](t, tc, engine)
			if tc.Assertion.StatusCode == http.StatusOK {
				require.Equal(t, claims.UserID, res.Body.Data.ForceMap("claims").ForceString("user_id"))
			}
		})
	}
}
