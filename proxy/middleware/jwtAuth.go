package middleware

import (
	"net/http"
	"strings"

	"github.com/erdsea/erdsea-api/crypto"
	"github.com/erdsea/erdsea-api/data"
	"github.com/gin-gonic/gin"
)

const (
	noBearerPresent = "No authorization bearer provided"
	incorrectBearer = "Incorrect bearer provided"
	invalidJwtToken = "Invalid or expired token"

	bearerSplitOn = "Bearer "
	authHeaderKey = "Authorization"

	addressKey = "address"
)

var returnUnauthorized = func(c *gin.Context, errMessage string) {
	data.JsonResponse(c, http.StatusUnauthorized, nil, errMessage)
}

func Authorization(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.Request.Header.Get(authHeaderKey)
		if bearer == "" {
			returnUnauthorized(c, noBearerPresent)
			c.Abort()
			return
		}

		ok, token := parseBearer(bearer)
		if !ok {
			returnUnauthorized(c, incorrectBearer)
			c.Abort()
			return
		}

		claims, err := crypto.ValidateJwt(token, secret)
		if err != nil {
			returnUnauthorized(c, invalidJwtToken)
			c.Abort()
			return
		}

		c.Set(addressKey, claims.Address)
		c.Next()
	}
}

func parseBearer(bearer string) (bool, string) {
	splitBearer := strings.Split(bearer, bearerSplitOn)

	if len(splitBearer) != 2 {
		return false, ""
	}

	return true, strings.TrimSpace(splitBearer[1])
}
