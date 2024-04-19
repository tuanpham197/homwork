package middleware

import (
	sctx "Ronin/component/appctx"
	mysqlRepoAuth "Ronin/internal/auth/infras/mysql"
	"Ronin/pkg/contants"
	"Ronin/pkg/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

func TokenVerificationMiddleware(ctx sctx.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		authorizationHeader := c.GetHeader("Authorization")

		// Check if the Authorization header is present
		if authorizationHeader == "" {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, constants.MissingTokenHeader, nil))
			c.Abort()
			return
		}

		// Extract the token from the Authorization header
		headerAuthorization := strings.Split(authorizationHeader, " ")
		if len(headerAuthorization) < 2 {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, constants.InvalidTokenHeader, nil))
		}

		tokenString := strings.TrimSpace(headerAuthorization[1])
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, constants.MissingToken, nil))
			c.Abort()
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			return []byte(os.Getenv("SECRET_KEY")), nil // Replace with your actual secret key
		})

		// Check for parsing or validation errors
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, constants.InvalidToken, err))
			c.Abort()
			return
		}

		// Check if the token is valid
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, constants.InvalidTokenClaim, nil))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, constants.InvalidTokenClaim, nil))
			c.Abort()
			return
		}

		userId, ok := claims["userId"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, constants.InvalidUserIDClaim, nil))
			c.Abort()
			return
		}
		// get info user
		repositoryAuth := mysqlRepoAuth.NewMySQLRepo(ctx.GetDBConnection(), ctx.GetRedisConnection())
		tokenUser, errGet := repositoryAuth.GetTokenUser(c.Request.Context(), map[string]interface{}{
			"user_id":    userId,
			"token_sign": strings.Split(tokenString, ".")[2],
		})
		if errGet != nil {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, "Asd", errGet))
			c.Abort()
			return
		}
		if tokenUser.Revoked == true {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, constants.InvalidToken, err))
			c.Abort()
			return
		}

		// Set info to redis

		roles, ok := claims["roles"]
		c.Set("roles", roles)

		c.Set("userId", userId)
		shopId := claims["shopId"]
		c.Set("shopId", shopId)

		// Call the next handler
		c.Next()
	}
}
