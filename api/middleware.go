package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/datewu/set-img/auth"
	"github.com/datewu/set-img/author"
	"github.com/gin-gonic/gin"
)

func checkAuth(c *gin.Context) {
	token, err := extractToken(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ok, err := auth.Valid(token)
	if err != nil || !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			gin.H{"token": token, "message": "bad token, cannot authentication"})
		return
	}
	ok, err = author.Can(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "has no authorization"})
		return
	}
	c.Next()
}

func extractToken(c *gin.Context) (string, error) {
	q := c.Query("token") // query
	if q != "" {
		return q, nil
	}
	ah := struct { // header
		Authorization string `header:"Authorization"`
	}{}
	err := c.ShouldBindHeader(&ah)
	if err != nil {
		return "", err
	}
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(ah.Authorization, bearerPrefix) {
		return "", errors.New("not Bearer token")
	}
	return ah.Authorization[len(bearerPrefix):], nil
}
