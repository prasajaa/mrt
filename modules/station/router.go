package station

import (
	"mrt-schedule/common/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitiateRouter(router *gin.RouterGroup) {
	StationService := NewService()
	station := router.Group("/stations")
	station.GET("/", func(c *gin.Context) {
		//code service
		GetAllStation(c, StationService)
	})

}

func GetAllStation(c *gin.Context, service Service) {
	data, err := service.GetAllStation()
	if err != nil {
		response := response.APIResponse("Error to get stations", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := response.APIResponse("List of stations", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
