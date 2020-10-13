package api

import (
	"net/http"
	"strings"

	"github.com/datewu/set-img/auth"
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
	api.GET("/token", getToken)
	private := api.Group("/auth", checkAuth)
	setPrivateRoutes(private)
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func getToken(c *gin.Context) {
	// TODO needs more security
	if !strings.HasPrefix(c.Request.Host, "localhost:") {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	token, err := auth.NewToken()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, token)
}

func setPrivateRoutes(api *gin.RouterGroup) {
	api.GET("/list/:ns", listDemo)
	api.POST("/setdeploy/:ns/image", setDeployImg)
}

func listDemo(c *gin.Context) {
	ns := c.Param("ns")
	ls, err := k8s.ListDemo(ns)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, ls)
}
func setDeployImg(c *gin.Context) {
	id := new(k8s.ContainerPath)
	if err := c.BindJSON(id); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"message": "cannot unmarshal/binding post data"})
		return
	}
	id.Ns = c.Param("ns")
	err := k8s.SetDeployImg(id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.String(http.StatusOK, "ok")
}
