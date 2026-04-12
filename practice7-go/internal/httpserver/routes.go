package httpserver

import (
	"net/http"
	"practice7/internal/middleware"
	"practice7/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterRoutes(r *gin.Engine, auth *utils.AuthService, jwtSecret []byte) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	users := r.Group("/users")
	{
		users.POST("", func(c *gin.Context) { register(c, auth) })
		users.POST("/login", func(c *gin.Context) { login(c, auth) })

		protected := users.Group("")
		protected.Use(middleware.JWTAuthMiddleware(jwtSecret))
		{
			protected.GET("/protected/hello", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "OK"})
			})

			// Task 1: GetMe
			protected.GET("/me", func(c *gin.Context) { getMe(c, auth) })

			// Task 2: Promote (admin-only)
			protected.PATCH("/promote/:id",
				middleware.RoleMiddleware("admin"),
				func(c *gin.Context) {
					idStr := c.Param("id")
					id, err := uuid.Parse(idStr)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
						return
					}
					user, err := auth.PromoteUserToAdmin(id)
					if err != nil {
						c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, gin.H{"user": user})
				},
			)
		}
	}
}

func register(c *gin.Context, auth *utils.AuthService) {
	var req utils.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := auth.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func login(c *gin.Context, auth *utils.AuthService) {
	var req utils.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := auth.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

func getMe(c *gin.Context, auth *utils.AuthService) {
	// Requirement: no input except JWT token
	userIDStr := c.GetString("userID")
	id, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id in token"})
		return
	}

	user, err := auth.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
