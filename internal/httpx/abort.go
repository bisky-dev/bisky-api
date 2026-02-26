package httpx

import (
	"github.com/gin-gonic/gin"
	"github.com/keithics/devops-dashboard/api/internal/httperr"
)

func AbortIfErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	httperr.Abort(c, httperr.BadRequest(err.Error()))
	return true
}
