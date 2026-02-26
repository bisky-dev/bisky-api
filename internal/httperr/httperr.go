package httperr

import (
	"errors"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func (e *HTTPError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

func (e *HTTPError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func New(status int, code string, message string) *HTTPError {
	return &HTTPError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func (e *HTTPError) WithDetails(details any) *HTTPError {
	e.Details = details
	return e
}

func (e *HTTPError) WithCause(err error) *HTTPError {
	e.Cause = err
	return e
}

func BadRequest(message string) *HTTPError {
	return New(http.StatusBadRequest, "BAD_REQUEST", message)
}

func Unauthorized(message string) *HTTPError {
	return New(http.StatusUnauthorized, "UNAUTHORIZED", message)
}

func NotFound(message string) *HTTPError {
	return New(http.StatusNotFound, "NOT_FOUND", message)
}

func Validation(message string) *HTTPError {
	return New(http.StatusUnprocessableEntity, "VALIDATION_ERROR", message)
}

func Conflict(message string) *HTTPError {
	return New(http.StatusConflict, "CONFLICT", message)
}

func PayloadTooLarge(message string) *HTTPError {
	return New(http.StatusRequestEntityTooLarge, "PAYLOAD_TOO_LARGE", message)
}

func Internal(message string) *HTTPError {
	return New(http.StatusInternalServerError, "INTERNAL_ERROR", message)
}

func From(err error) *HTTPError {
	if err == nil {
		return nil
	}
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return httpErr
	}
	return Internal("Internal Server Error").WithCause(err)
}

func Write(c *gin.Context, err error) {
	httpErr := From(err)
	if httpErr == nil {
		return
	}
	c.JSON(httpErr.Status, APIErrorResponse{
		Error: APIError{
			Code:    httpErr.Code,
			Message: httpErr.Message,
			Details: httpErr.Details,
		},
	})
}

func Abort(c *gin.Context, err error) {
	_ = c.Error(err)
	c.Abort()
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Written() || len(c.Errors) == 0 {
			return
		}
		Write(c, c.Errors.Last().Err)
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic recovered: %v\n%s", rec, string(debug.Stack()))
				Abort(c, Internal("Internal Server Error"))
				if !c.Writer.Written() {
					Write(c, c.Errors.Last().Err)
				}
			}
		}()
		c.Next()
	}
}
