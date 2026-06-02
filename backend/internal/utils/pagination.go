package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPage(c *gin.Context) int {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	return page
}

func GetPageSize(c *gin.Context) int {
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return size
}

func GetOffset(page, pageSize int) int {
	return (page - 1) * pageSize
}
