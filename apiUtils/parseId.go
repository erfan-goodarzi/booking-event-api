package apiUtils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func ParseID(c *gin.Context) (string, error) {
	idParams := c.Param("id")

	if idParams == "" {
		return "", errors.New("Invalid id")
	}

	return idParams, nil
}
