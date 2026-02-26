package httpx

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/keithics/devops-dashboard/api/internal/httperr"
)

func AbortDBErr(c *gin.Context, err error, internalMsg string) {
	if errors.Is(err, pgx.ErrNoRows) {
		httperr.Abort(c, httperr.NotFound("resource not found").WithCause(err))
		return
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			httperr.Abort(c, httperr.Conflict("resource already exists").WithCause(err))
			return
		case "23503", "23514":
			httperr.Abort(c, httperr.BadRequest("invalid input").WithCause(err))
			return
		}
	}

	httperr.Abort(c, httperr.Internal(internalMsg).WithCause(err))
}

func AbortIfDBErr(c *gin.Context, err error, internalMsg string) bool {
	if err == nil {
		return false
	}

	AbortDBErr(c, err, internalMsg)
	return true
}

func AbortDBErrNotFoundMsg(c *gin.Context, err error, notFoundMsg, internalMsg string) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, pgx.ErrNoRows) {
		httperr.Abort(c, httperr.NotFound(notFoundMsg).WithCause(err))
		return true
	}

	AbortDBErr(c, err, internalMsg)
	return true
}
