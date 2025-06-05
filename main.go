package main

import (
	"mrt-schedule/modules/station"

	"github.com/gin-gonic/gin"
)

func main() {
	InitiateRouter()
}

func InitiateRouter() {

	r := gin.Default()
	api := r.Group("/api/v1")

	station.InitiateRouter(api)

	r.Run()
}
