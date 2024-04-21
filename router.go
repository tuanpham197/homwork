package main

import (
	"Ronin/component/appctx"
	apphasher "Ronin/component/hasher"
	"Ronin/internal/booking/controller/httpapi"
	middlewareauth "Ronin/internal/booking/controller/middleware"
	mysqlbooking "Ronin/internal/booking/infras/mysql"
	serviceBooking "Ronin/internal/booking/service"
	"Ronin/middleware"
	"github.com/gin-gonic/gin"

	apiauth "Ronin/internal/auth/controller/api"
	mysqlauth "Ronin/internal/auth/infras/mysql"
	serviceauth "Ronin/internal/auth/services"
)

func setupBookingRoute(appCtx appctx.AppCtx, v1 *gin.RouterGroup) {

	repo := mysqlbooking.NewMySQLRepo(appCtx.GetMainDBConnection())
	biz := serviceBooking.NewService(repo)
	api := httpapi.NewAPIController(biz)

	v1.Use(middlewareauth.TokenVerificationMiddleware(appCtx), middleware.RoleMiddleware(appCtx, "admin"))
	v1.POST("/bookings", api.CreateBooking())
}

func setupAuthRoute(appCtx appctx.AppCtx, v1 *gin.RouterGroup) {
	repo := mysqlauth.NewMySQLRepo(appCtx.GetMainDBConnection(), appCtx.GetRedisClient())
	hasher := new(apphasher.Hasher)
	biz := serviceauth.NewService(repo, hasher, nil, appCtx.GetRedisClient())

	api := apiauth.NewAPIController(biz)

	v1.POST("/login", api.Login())
}
