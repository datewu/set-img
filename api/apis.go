package api

import (
	"net/http"

	"github.com/datewu/set-img/k8s"
	"github.com/gin-gonic/gin"
)

func getEngine() *gin.Engine {
	r := gin.Default()

	api := r.Group(apiVersion)
	setRoutes(api)
	return r
}

func setRoutes(api *gin.RouterGroup) {
	api.GET("/ping", ping)
	api.GET("/list/:ns", listDemo)
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func listDemo(c *gin.Context) {
	ns := c.Param("ns")
	c.JSON(http.StatusOK, k8s.ListDemo(ns))
}
