package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	apiVersion  = "/api/v1"
	maxBodySize = 10 << 20 // 10mb
)

// Conf api configuration
type Conf struct {
	Addr string
	Mode string
}

// Server ...
func Server(c *Conf) error {
	if c.Mode == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	e := getEngine()
	e.MaxMultipartMemory = maxBodySize
	if c.Mode == "dev" {
		corsMiddler := func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			c.Next()
		}
		e.Use(corsMiddler)
	}
	s := &http.Server{
		Addr:           c.Addr,
		Handler:        e,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   35 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return s.ListenAndServe()
}
