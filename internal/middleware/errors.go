package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	error2 "gitlab.unanet.io/devops/eve/internal/error"
	"net/http"
)

func ApiError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}

		var restError *error2.RestError
		if errors.As(err, &restError) {
			// TODO: ErrorHandling, Under Debug we should also log the original Error that this is wrapping
			c.AbortWithStatusJSON(restError.Code, restError)
			return
		}

		// TODO: ErrorHandling Log Error, something unexpected happened
		c.AbortWithStatusJSON(http.StatusInternalServerError, error2.RestError{Code: 500, Message: "Internal Server Error"})

	}
}