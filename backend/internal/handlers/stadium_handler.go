package handlers

import (
	"strconv"

	"worldcup-mate/internal/repositories"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
)

func ListStadiums(c *gin.Context) {
	stadiums, err := repositories.ListStadiums()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, stadiums)
}

func GetStadiumDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid stadium id")
		return
	}
	stadium, err := repositories.GetStadiumByID(uint(id))
	if err != nil {
		utils.Error(c, 404, "stadium not found")
		return
	}
	utils.Success(c, stadium)
}
