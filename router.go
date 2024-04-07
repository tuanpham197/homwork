package main

import (
	"Ronin/component/appctx"
	"Ronin/modules/booking/controller/httpapi"
	mysqlbooking "Ronin/modules/booking/infras/mysql"
	serviceBooking "Ronin/modules/booking/service"
	"github.com/gin-gonic/gin"
)

func setupBookingRoute(appCtx appctx.AppCtx, v1 *gin.RouterGroup) {

	repo := mysqlbooking.NewMySQLRepo(appCtx.GetMainDBConnection())
	biz := serviceBooking.NewService(repo)
	api := httpapi.NewAPIController(biz)

	v1.POST("/bookings", api.CreateBooking())
}
