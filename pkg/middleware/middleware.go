package middleware

import (
    "freepass-2026/internal/service"
    "freepass-2026/pkg/jwt"

    "github.com/gin-gonic/gin"
)

type Interface interface {
    AuthenticateUser(c *gin.Context)
    AuthenticateAdmin(c *gin.Context)
    AuthenticateCanteenOwner(c *gin.Context)
}

type middleware struct {
    service *service.Service
    jwtAuth jwt.Interface
}

func Init(service *service.Service, jwtAuth jwt.Interface) Interface {
    return &middleware{
        service: service,
        jwtAuth: jwtAuth,
    }
}