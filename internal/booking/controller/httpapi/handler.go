package httpapi

import (
	"Ronin/internal/booking/entity"
	"Ronin/internal/booking/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type apiController struct {
	service service.BookingUseCase
}

func NewAPIController(s service.BookingUseCase) *apiController {
	return &apiController{service: s}
}

func (api apiController) CreateBooking() func(c *gin.Context) {

	return func(c *gin.Context) {
		var request entity.BookingCreate

		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		booking, err := api.service.CreateBooking(c.Request.Context(), &request)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, booking)
	}
}
