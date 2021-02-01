package token

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/tianhongw/tinyid/pkg/errdef"
)

type TokenService interface {
	GetTokenByBizType(bizType string) (token string, err error)
}

func Authentication(tokenService TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenFromQuery := c.Query("token")
		if tokenFromQuery == "" {
			c.Error(errdef.HttpErrorUnauthorized).SetType(gin.ErrorTypePublic)
			c.Abort()
			return
		}

		bizType := c.Query("type")
		tokenFromDb, err := tokenService.GetTokenByBizType(bizType)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePublic)
			c.Abort()
			return
		}

		if tokenFromQuery != tokenFromDb {
			c.Error(errors.New("token not match")).SetType(gin.ErrorTypePublic)
			c.Abort()
			return
		}

		c.Next()
	}
}
