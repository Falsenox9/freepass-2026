package rest

import (
    "fmt"
    "freepass-2026/internal/service"
    "freepass-2026/pkg/middleware"
    "os"

    "github.com/gin-gonic/gin"
)

type Rest struct {
    router     *gin.Engine
    service    *service.Service
    middleware middleware.Interface
}

func NewRest(service *service.Service, middleware middleware.Interface) *Rest {
    return &Rest{
        router:     gin.Default(),
        service:    service,
        middleware: middleware,
    }
}

func (r *Rest) MountEndPoint() {
    baseURL := r.router.Group("/api/v1")

    auth := baseURL.Group("/auth")
    {
        auth.POST("/register", r.RegisterUser)
        auth.POST("/login", r.LoginUser)
    }

    user := baseURL.Group("/user")
    user.Use(r.middleware.AuthenticateUser)
    {
        // Add user-specific routes here if needed
        // user.GET("/profile", r.GetUserProfile)
    }

    admin := baseURL.Group("/admin")
    admin.Use(r.middleware.AuthenticateAdmin)
    {
        admin.POST("/canteen-owner", r.CreateCanteenOwner)
        admin.PUT("/canteen-owner", r.UpdateCanteenOwner)
        admin.DELETE("/user/:user_id", r.DeleteUser)
        admin.GET("/users", r.GetAllUsers)
    }

    canteen := baseURL.Group("/canteen")
    canteen.Use(r.middleware.AuthenticateCanteenOwner)
    {
        canteen.POST("/food", r.CreateFood)
        canteen.PUT("/food", r.UpdateFood)
        canteen.DELETE("/food/:food_id", r.DeleteFood)
        canteen.GET("/my-foods", r.GetMyFoods)
        canteen.PATCH("/food/stock", r.UpdateStock)
    }

    public := baseURL.Group("/public")
    {
        public.GET("/foods", r.GetAllFoods)
        public.GET("/food/:food_id", r.GetFoodByID)
    }
}

func (r *Rest) Run() {
    addr := os.Getenv("ADDRESS")
    port := os.Getenv("PORT")

    r.router.Run(fmt.Sprintf("%s:%s", addr, port))
}