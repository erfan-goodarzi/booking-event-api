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

func ParsParam(c *gin.Context, name string) (string, error) {
	param := c.Param(name)

	if param == "" {
		return "", errors.New("Invalid" + name)
	}

	return param, nil
}
