package middleware

import (
	"net/http"
	"strings"

	"github.com/ENFT-DAO/youbei-api/crypto"
	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/services"
	"github.com/gin-gonic/gin"
)

const (
	noBearerPresent = "No authorization bearer provided"
	incorrectBearer = "Incorrect bearer provided"
	invalidJwtToken = "Invalid or expired token"

	bearerSplitOn = "Bearer "
	authHeaderKey = "Authorization"

	AddressKey = "address"
	IsAdminKey = "isAdmin"
)

var returnUnauthorized = func(c *gin.Context, errMessage string) {
	dtos.JsonResponse(c, http.StatusUnauthorized, nil, errMessage)
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

		//claims.Address
		// Get the account base on web-wallet address
		// check the role if it's "RoleAdmin"

		account, _ := services.GetOrCreateAccount(claims.Address)
		isRoleAdmin := (account.Role == "RoleAdmin")

		c.Set(AddressKey, claims.Address)
		c.Set(IsAdminKey, isRoleAdmin)
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
