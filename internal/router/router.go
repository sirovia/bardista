package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sirovia/bardista/internal/handler"
)

func Setup(
	r *gin.Engine,
	authHandler *handler.AuthHandler,
	productHandler *handler.ProductHandler,
	jwtSecret string,
) {
	public := r.Group("/api/v1")
	authed := r.Group("/api/v1", handler.AuthMiddleware(jwtSecret))
	adminOnly := r.Group("/api/v1", handler.AuthMiddleware(jwtSecret), handler.AdminMiddleware())

	_ = authed

	public.POST("/auth/register", authHandler.Register)
	public.POST("/auth/login", authHandler.Login)

	public.GET("/products", productHandler.GetAll)
	public.GET("/products/:id", productHandler.GetByID)
	adminOnly.POST("/products", productHandler.Create)
	adminOnly.PUT("/products/:id", productHandler.Update)
	adminOnly.DELETE("/products/:id", productHandler.Delete)
}
