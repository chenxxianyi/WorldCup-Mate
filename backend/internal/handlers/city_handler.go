package handlers

import (
	"worldcup-mate/internal/repositories"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
)

func ListCities(c *gin.Context) {
	cities, err := repositories.ListCities()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, cities)
}
