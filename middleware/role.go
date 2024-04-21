package middleware

import (
	commonapp "Ronin/common"
	sctx "Ronin/component/appctx"
	"Ronin/internal/auth/services/entity"
	"Ronin/pkg/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"net/http"
	"strings"
)

func RoleMiddleware(ctx sctx.AppCtx, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// user info
		userId := c.GetString("userId")

		if userId == "" {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, "Unauthorized", nil))
			c.Abort()
			return
		}

		// Get info user
		key := fmt.Sprintf("user_info:%s", userId)
		var userInfo entity.User
		if err := commonapp.GetDatRedis(c, key, &userInfo, ctx.GetRedisClient()); err != nil {
			c.JSON(http.StatusForbidden, response.ErrorResponse(http.StatusForbidden, "Forbidden", nil))
			c.Abort()
			return
		}

		roles := lo.FlatMap(userInfo.Roles, func(item entity.Role, index int) []string {
			return []string{strings.ToLower(item.Name)}
		})

		fmt.Println(roles, "abc")

		if !commonapp.CheckExistsRole(roles, role) {
			c.JSON(http.StatusForbidden, response.ErrorResponse(http.StatusForbidden, "Forbidden", nil))
			c.Abort()
			return
		}

		c.Next()
	}
}
