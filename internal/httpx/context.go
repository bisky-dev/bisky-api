package httpx

import (
	"github.com/gin-gonic/gin"
	"github.com/keithics/devops-dashboard/api/internal/httperr"
)

func Get[T any](c *gin.Context, key string) (T, bool) {
	v, ok := c.Get(key)
	if !ok {
		var zero T
		return zero, false
	}

	val, ok := v.(T)
	if !ok {
		var zero T
		return zero, false
	}

	return val, true
}

func MustGet[T any](c *gin.Context, key string) T {
	v, ok := c.Get(key)
	if !ok {
		panic("missing context value: " + key)
	}

	val, ok := v.(T)
	if !ok {
		panic("invalid context value type: " + key)
	}

	return val
}

func AbortIfMissingContext[T any](c *gin.Context, key string) (T, bool) {
	val, ok := Get[T](c, key)
	if ok {
		return val, true
	}

	httperr.Abort(c, httperr.Internal("missing request context"))
	var zero T
	return zero, false
}
