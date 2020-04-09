package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gitlab.unanet.io/devops/eve/internal/common"
	"net/http"
)

func ApiError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}

		var restError *common.RestError
		if errors.As(err, &restError) {
			// TODO: ErrorHandling, Under Debug we should also log the original Error that this is wrapping
			c.AbortWithStatusJSON(restError.Code, restError)
			return
		}

		// TODO: ErrorHandling Log Error, something unexpected happened
		c.AbortWithStatusJSON(http.StatusInternalServerError, common.RestError{Code: 500, Message: "Internal Server Error"})

	}
}