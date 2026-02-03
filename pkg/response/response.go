package response

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalItems int64       `json:"total_items"`
	TotalPages int         `json:"total_pages"`
}

func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, errMessage string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   errMessage,
	})
}

func Paginated(c *gin.Context, statusCode int, data interface{}, page, limit int, totalItems int64) {
	totalPages := int(totalItems) / limit
	if int(totalItems)%limit > 0 {
		totalPages++
	}

	c.JSON(statusCode, PaginatedResponse{
		Success:    true,
		Data:       data,
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	})
}
