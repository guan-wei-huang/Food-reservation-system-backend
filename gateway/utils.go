package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func getUserId(c *gin.Context) int {
	uid := c.GetString("userId")
	id, _ := strconv.Atoi(uid)
	return id
}
