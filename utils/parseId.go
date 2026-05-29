package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseID(c *gin.Context) (int64, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
