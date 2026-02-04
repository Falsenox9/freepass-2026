package middleware

import (
    "errors"
    "freepass-2026/model"
    "freepass-2026/pkg/response"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

func (m *middleware) AuthenticateAdmin(c *gin.Context) {
    bearer := c.GetHeader("Authorization")
    if bearer == "" {
        response.Error(c, http.StatusUnauthorized, "empty token", nil)
        c.Abort()
        return
    }

    token := strings.Split(bearer, " ")
    if len(token) != 2 {
        response.Error(c, http.StatusUnauthorized, "invalid token format", nil)
        c.Abort()
        return
    }

    userID, err := m.jwtAuth.ValidateToken(token[1])
    if err != nil {
        response.Error(c, http.StatusUnauthorized, "invalid token", nil)
        c.Abort()
        return
    }

    user, err := m.service.UserService.GetUser(model.UserParam{
        UserID: userID,
    })

    if err != nil {
        response.Error(c, http.StatusUnauthorized, "failed to get user", nil)
        c.Abort()
        return
    }

    if user.RoleID != 1 {
        response.Error(c, http.StatusForbidden, "this endpoint cant be access", errors.New("user dont have access"))
        c.Abort()
        return
    }

    c.Set("user", user)
    c.Set("user_id", user.UserID)
    c.Next()
}

func (m *middleware) AuthenticateCanteenOwner(c *gin.Context) {
    bearer := c.GetHeader("Authorization")
    if bearer == "" {
        response.Error(c, http.StatusUnauthorized, "empty token", nil)
        c.Abort()
        return
    }

    token := strings.Split(bearer, " ")
    if len(token) != 2 {
        response.Error(c, http.StatusUnauthorized, "invalid token format", nil)
        c.Abort()
        return
    }

    userID, err := m.jwtAuth.ValidateToken(token[1])
    if err != nil {
        response.Error(c, http.StatusUnauthorized, "invalid token", nil)
        c.Abort()
        return
    }

    user, err := m.service.UserService.GetUser(model.UserParam{
        UserID: userID,
    })

    if err != nil {
        response.Error(c, http.StatusUnauthorized, "failed to get user", nil)
        c.Abort()
        return
    }

    if user.RoleID != 3 {
        response.Error(c, http.StatusForbidden, "this endpoint cant be access", errors.New("user dont have access"))
        c.Abort()
        return
    }

    c.Set("user", user)
    c.Set("user_id", user.UserID)
    c.Next()
}