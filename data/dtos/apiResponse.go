package dtos

import "github.com/gin-gonic/gin"

type ApiResponse struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

func JsonResponse(c *gin.Context, status int, data interface{}, error string) {
	c.JSON(status, ApiResponse{
		Data:  data,
		Error: error,
	})
}
