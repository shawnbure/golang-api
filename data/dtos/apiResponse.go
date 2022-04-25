package dtos

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

func ContentAsFileResponse(c *gin.Context, filename string, data *bytes.Buffer) {
	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprintf("%d", len(data.Bytes())))
	c.Writer.Write(data.Bytes())
}

func StringResponse(c *gin.Context, data string) {
	c.String(http.StatusOK, "%v", data)
}
