package response

import (
    "github.com/gin-gonic/gin"
)

type Status struct {
    Code      int  `json:"code"`
    IsSuccess bool `json:"isSuccess"`
}

type Response struct {
    Status  Status      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

func Success(ctx *gin.Context, code int, message string, data any) {
    ctx.JSON(code, Response{
        Status: Status{
            Code:      code,
            IsSuccess: true,
        },
        Message: message,
        Data:    data,
    })
}

func Error(ctx *gin.Context, code int, message string, err error) {
    var errorData interface{}
    if err != nil {
        errorData = err.Error()
    } else {
        errorData = nil
    }

    ctx.JSON(code, Response{
        Status: Status{
            Code:      code,
            IsSuccess: false,
        },
        Message: message,
        Data:    errorData,
    })
}