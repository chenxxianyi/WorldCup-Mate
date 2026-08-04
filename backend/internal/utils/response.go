package utils

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	RequestID string      `json:"request_id,omitempty"`
}

type PaginatedData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Error writes a business error. API-01: the HTTP status code equals the
// business code (400/401/403/404/409/429/500...). Server-side errors (5xx)
// are logged with their internal detail and the client only receives a
// generic message — internal details never leak to the client.
func Error(c *gin.Context, code int, message string) {
	if code >= 500 {
		log.Printf("[api-error] status=%d detail=%s", code, message)
		message = "internal server error"
	}
	requestID, _ := c.Get("request_id")
	c.JSON(code, Response{
		Code:      code,
		Message:   message,
		Data:      nil,
		RequestID: fmt.Sprintf("%v", requestID),
	})
}

func Paginated(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	Success(c, PaginatedData{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}
