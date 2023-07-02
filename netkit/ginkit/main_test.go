package ginkit

import (
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	Setup()
	goleak.VerifyTestMain(m)
}
