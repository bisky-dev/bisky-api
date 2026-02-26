package httpx

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/mold/v4"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/validator/v10"
	"github.com/keithics/devops-dashboard/api/internal/httperr"
)

var (
	validate   = newValidator()
	normalizer = newNormalizer()
)

func newValidator() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := field.Tag.Get("json")
		if name == "" {
			return field.Name
		}

		name = strings.Split(name, ",")[0]
		if name == "" || name == "-" {
			return field.Name
		}

		return name
	})

	return v
}

func newNormalizer() *mold.Transformer {
	return modifiers.New()
}

func NormalizeAndValidate(input any) error {
	if err := normalizer.Struct(context.Background(), input); err != nil {
		return err
	}
	if err := validate.Struct(input); err != nil {
		return err
	}
	return nil
}

func AbortIfNormalizeValidateErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	var validationErr validator.ValidationErrors
	if errors.As(err, &validationErr) {
		httperr.Abort(c, httperr.Validation("validation failed").WithCause(err))
		return true
	}

	httperr.Abort(c, httperr.BadRequest("invalid input").WithCause(err))
	return true
}
