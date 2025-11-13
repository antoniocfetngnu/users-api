package main

import (
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/antoniocfetngnu/users-api/config"
	"github.com/antoniocfetngnu/users-api/database"
	_ "github.com/antoniocfetngnu/users-api/docs" // Swagger docs
	"github.com/antoniocfetngnu/users-api/graphql"
	"github.com/antoniocfetngnu/users-api/handlers"
	"github.com/antoniocfetngnu/users-api/utils"
)

// @title Users Service API
// @version 1.0
// @description Users microservice with REST API and GraphQL
// @host localhost:3001
// @BasePath /
// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name auth_token
func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize JWT utils
	utils.InitJWT(cfg)

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Setup Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:8000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Middleware to extract user from cookie (for protected routes)
	authMiddleware := func(c *gin.Context) {
		cookie, err := c.Cookie("auth_token")
		if err != nil || cookie == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateJWT(cookie)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Next()
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"service": "users-service",
		})
	})

	// Public routes (no authentication)
	r.POST("/api/auth/register", handlers.Register)
	r.POST("/api/auth/login", handlers.Login)
	r.POST("/api/auth/logout", handlers.Logout)

	// Protected REST routes (require authentication)
	authorized := r.Group("/api/users")
	authorized.Use(authMiddleware)
	{
		authorized.GET("", handlers.GetUsers)
		authorized.GET("/:id", handlers.GetUser)
		authorized.PUT("/:id", handlers.UpdateUser)
		authorized.DELETE("/:id", handlers.DeleteUser)
	}

	// GraphQL setup
	gqlResolver := &graphql.Resolver{}
	gqlServer := handler.NewDefaultServer(
		graphql.NewExecutableSchema(graphql.Config{Resolvers: gqlResolver}),
	)

	// GraphQL endpoint (protected)
	r.POST("/graphql", authMiddleware, func(c *gin.Context) {
		gqlServer.ServeHTTP(c.Writer, c.Request)
	})

	// GraphQL playground (for development)
	if cfg.Environment == "development" {
		r.GET("/playground", func(c *gin.Context) {
			playground.Handler("GraphQL Playground", "/graphql").ServeHTTP(c.Writer, c.Request)
		})
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("🚀 Server running on http://localhost:%s", cfg.Port)
	log.Printf("📊 Health check: http://localhost:%s/health", cfg.Port)
	log.Printf("📚 Swagger UI: http://localhost:%s/swagger/index.html", cfg.Port)
	if cfg.Environment == "development" {
		log.Printf("🎮 GraphQL Playground: http://localhost:%s/playground", cfg.Port)
	}

	r.Run(":" + cfg.Port)
}
