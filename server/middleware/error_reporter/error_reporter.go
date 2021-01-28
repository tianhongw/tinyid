package error_reporter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tianhongw/tinyid/pkg/errdef"
)

/**
 * Auto response error
 *
 * Usage:
 * c.Error(err).SetType(...)
 * c.Abort()
 */
func JSONErrorReporter() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		lastErr := c.Errors.Last()

		// do not handle when response already written
		if lastErr == nil || c.Writer.Written() {
			return
		}

		var httpError errdef.HttpError
		var status int

		switch parsedErr := lastErr.Err.(type) {
		case errdef.HttpError:
			status = parsedErr.Code.StatusCode()
			httpError = parsedErr
		default:
			status = http.StatusInternalServerError
			httpError = errdef.NewHttpError(parsedErr.Error(), errdef.HttpErrorCodeUniversal)
		}

		if c.Writer.Status() != http.StatusOK {
			status = c.Writer.Status()
		}

		c.JSON(status, httpError)
	}
}
