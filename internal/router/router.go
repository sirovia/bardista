package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sirovia/bardista/internal/handler"
)

func Setup(
	r *gin.Engine,
	authHandler *handler.AuthHandler,
	jwtSecret string,
) {
	public := r.Group("/api/v1")
	authed := r.Group("/api/v1", handler.AuthMiddleware(jwtSecret))
	adminOnly := r.Group("/api/v1", handler.AuthMiddleware(jwtSecret), handler.AdminMiddleware())

	_ = authed
	_ = adminOnly

	public.POST("/auth/register", authHandler.Register)
	public.POST("/auth/login", authHandler.Login)
}
