package api

import (
	service "Ronin/internal/auth/services"
	"Ronin/internal/auth/services/entity"
	"Ronin/internal/auth/services/request"
	"Ronin/pkg/utils/response"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type apiController struct {
	service service.UserUseCase
}

func NewAPIController(s service.UserUseCase) apiController {
	return apiController{service: s}
}

func (api apiController) Login() func(ctx *gin.Context) {

	return func(c *gin.Context) {
		var userLogin request.UserLogin
		errBind := c.ShouldBindJSON(&userLogin)
		if errBind != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": errBind.Error(),
			})
			return
		}
		result, err := api.service.Login(c, userLogin)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, response.ResponseData(result, http.StatusOK, ""))
	}
}

func (api apiController) Register() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var userRegister request.UserRegister

		errBind := c.ShouldBindJSON(&userRegister)
		if errBind != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": errBind,
			})
			return
		}
		result, err := api.service.Register(c, userRegister)
		if err != nil {
			var userExistsError request.UserExistsError
			if errors.As(err, &userExistsError) {
				c.JSON(http.StatusBadRequest, gin.H{
					"err": err.Error(),
				})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, response.ResponseData(result, http.StatusOK, ""))
	}
}

func (api apiController) GetInfoUser() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		userId := c.GetString("userId")

		user, err := api.service.GetInfoUser(c, userId)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, response.ResponseData(user, http.StatusOK, ""))
	}

}

func (api apiController) RefreshToken() func(c *gin.Context) {
	return func(c *gin.Context) {
		token, _ := c.GetPostForm("refresh_token")

		tokenRequest := request.RefreshTokenRequest{
			RefreshToken: token,
		}

		result, err := api.service.RefreshToken(c, tokenRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, response.ResponseData(result, http.StatusOK, ""))
	}
}

func (api apiController) RevokeToken() func(c *gin.Context) {
	return func(c *gin.Context) {
		var req entity.TokenUserUpdate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		userId := c.GetString("userId")
		req.UserId = userId
		if errUpdate := api.service.RevokeToken(c, &req); errUpdate != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": errUpdate.Error()})
			return
		}
		c.JSON(http.StatusOK, response.ResponseData(true, http.StatusOK, ""))
	}
}
