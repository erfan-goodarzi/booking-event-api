package apiUtils

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseID(c *gin.Context) (int64, error) {
	idParams := c.Param("id")

	if idParams == "" {
		return 0, errors.New("Invalid id")
	}

	id, err := strconv.ParseInt(idParams, 10, 64)

	if err != nil {
		return 0, err
	}

	return id, nil
}
