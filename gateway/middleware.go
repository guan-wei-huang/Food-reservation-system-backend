package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth(c *gin.Context) {
	token := c.Request.Header.Get("access_token")
	// refreshToken := c.Request.Header.Get("refresh_token")

	tokenClaims, err := jwt.ParseWithClaims(token, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	claims := tokenClaims.Claims.(*jwt.MapClaims)
	// exp := claims.VerifyExpiresAt(time.Now().Unix(), true)
	// if exp {
	// 	// generate new token
	// 	token = refreshToken
	// }

	uid := (*claims)["userId"]
	c.Set("userId", uid)
}
