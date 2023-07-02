package ginkit

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hungdv136/gokit/netkit"
	"github.com/hungdv136/gokit/netkit/testkit"
	"github.com/hungdv136/gokit/types"
	"github.com/stretchr/testify/require"
)

type (
	nested struct {
		Email string `json:"email,omitempty" binding:"required"`
	}

	payload struct {
		PhoneNumber string    `json:"phone_number,omitempty" binding:"required"`
		ID          int64     `json:"id" binding:"required"`
		Nested      nested    `json:"nested,omitempty" binding:"required,email" message:"test_message"`
		StartedAt   time.Time `json:"started_at" binding:"required"`
	}

	invalid struct {
		InvalidParameters []string `json:"invalid_parameters"`
	}
)

func TestValidation(t *testing.T) {
	t.Parallel()

	handler := func(ctx *gin.Context) {
		payload := payload{}
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			SendValidationError(ctx, err)
			return
		}

		SendSuccess(ctx, "success", payload)
	}

	engine := gin.New()
	engine.POST("/test", handler)

	validParam := types.Map{
		"phone_number": "09887788999",
		"id":           123,
		"nested": types.Map{
			"email": "support@ts.com",
		},
		"started_at": time.Now().Format(time.RFC3339),
	}

	invalidParameter := types.Map{
		"phone_number": "09887788999",
		"id":           "123",
		"started_at":   time.Now().Format(time.RFC3339),
	}

	invalidDateParam := types.Map{
		"phone_number": "09887788999",
		"id":           123,
		"nested": types.Map{
			"email": "support@ts.com",
		},
		"started_at": "2022",
	}

	testCases := []*testkit.TestCase{
		testkit.NewTestCase("success", http.MethodPost, "/test", http.StatusOK, netkit.VerdictSuccess).WithBody(validParam),
		testkit.NewTestCase("missing", http.MethodPost, "/test", http.StatusBadRequest, netkit.VerdictInvalidParameters).WithBody(types.Map{}),
		testkit.NewTestCase("invalid", http.MethodPost, "/test", http.StatusBadRequest, netkit.VerdictInvalidParameters).WithBody(invalidParameter),
		testkit.NewTestCase("invalid_date", http.MethodPost, "/test", http.StatusBadRequest, netkit.VerdictInvalidParameters).WithBody(invalidDateParam),
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			res := testkit.TestGin[invalid](t, tc, engine)
			fields := res.Body.Data.InvalidParameters

			if tc.Name == "missing" {
				require.ElementsMatch(t, []string{"id", "phone_number", "nested.email", "started_at"}, fields)
			} else if tc.Name == "invalid" {
				require.ElementsMatch(t, []string{"id"}, fields)
			} else if tc.Name != "invalid_date" {
				require.Empty(t, fields)
			}
		})
	}
}
