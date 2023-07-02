package ginkit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"hungdv136/gokit/logger"
	"hungdv136/gokit/netkit"
	"hungdv136/gokit/types"
)

const fieldMsg = "invalid '%s' tag"

type Field struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type Fields []*Field

func (f Fields) GetNames() []string {
	names := make([]string, len(f))
	for i, field := range f {
		names[i] = field.Name
	}

	return names
}

func (f Fields) ToKeyValues() types.Map {
	kv := make(types.Map, len(f))
	for _, field := range f {
		kv[field.Name] = field.Message
	}

	return kv
}

// RegisterTagNameFunc registers tags function
func RegisterTagNameFunc() bool {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return false
	}

	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})

	return true
}

// SendValidationError sends validation error response to client
func SendValidationError(ctx *gin.Context, err error) {
	f := GetInvalidParameters(err)
	if len(f) == 0 {
		logger.Info(ctx, "unknown validation error", err)
		SendJSON(ctx, http.StatusBadRequest, netkit.VerdictInvalidParameters, "invalid parameter", types.Map{})
		return
	}

	logger.Fields(ctx, "fields", f.ToKeyValues()).Info(ctx, "invalid parameters")
	SendInvalidParameters(ctx, f.GetNames())
}

// GetInvalidParameters gets list of invalid parameters
func GetInvalidParameters(err error) Fields {
	if field, ok := GetUnmarshalErrorParameters(err); ok {
		return []*Field{field}
	}

	vErrs := validator.ValidationErrors{}
	if !errors.As(err, &vErrs) {
		return nil
	}

	fields := make([]*Field, len(vErrs))
	for i, vErr := range vErrs {
		msg := fmt.Sprintf(fieldMsg, vErr.Tag())
		ns := vErr.Namespace()
		if len(ns) == 0 {
			fields[i] = &Field{Name: vErr.Field(), Message: msg}
			continue
		}

		if index := strings.Index(ns, "."); index > 0 {
			fields[i] = &Field{Name: ns[index+1:], Message: msg}
			continue
		}

		fields[i] = &Field{Name: ns, Message: msg}
	}

	return fields
}

// GetUnmarshalErrorParameters returns list of invalid parameters
func GetUnmarshalErrorParameters(err error) (*Field, bool) {
	var uErr *json.UnmarshalTypeError
	if errors.As(err, &uErr) {
		return &Field{Name: uErr.Field, Message: uErr.Error()}, true
	}

	return nil, false
}
