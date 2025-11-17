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
	_ "github.com/antoniocfetngnu/users-api/docs"
	"github.com/antoniocfetngnu/users-api/graphql"
	"github.com/antoniocfetngnu/users-api/handlers"
	"github.com/antoniocfetngnu/users-api/middleware"
	"github.com/antoniocfetngnu/users-api/utils"
)

// @title Users Service API
// @version 1.0
// @description Users microservice with REST API and GraphQL
// @host localhost:8000
// @BasePath /
// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name auth_token
func main() {
	cfg := config.LoadConfig()
	utils.InitJWT(cfg)

	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

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

	// Health check (public)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"service": "users-service",
		})
	})

	// Public auth routes (no authentication required)
	r.POST("/api/auth/register", handlers.Register)
	r.POST("/api/auth/login", handlers.Login)
	r.POST("/api/auth/logout", handlers.Logout)

	// Protected auth routes (require authentication)
	authProtected := r.Group("/api/auth")
	authProtected.Use(middleware.AuthMiddleware())
	{
		authProtected.GET("/me", handlers.Me)
	}

	// Protected user routes
	authorized := r.Group("/api/users")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.GET("", handlers.GetUsers)
		authorized.GET("/:id", handlers.GetUser)
		authorized.PUT("/:id", handlers.UpdateUser)
		authorized.DELETE("/:id", handlers.DeleteUser)
	}

	// Follower routes (protected)
	followers := r.Group("/api/followers")
	followers.Use(middleware.AuthMiddleware())
	{
		followers.POST("/follow", handlers.FollowUser)
		followers.DELETE("/unfollow/:id", handlers.UnfollowUser)
		followers.GET("/my-followers", handlers.GetMyFollowers)
		followers.GET("/my-following", handlers.GetMyFollowing)
	}

	// GraphQL setup
	gqlResolver := &graphql.Resolver{}
	gqlServer := handler.NewDefaultServer(
		graphql.NewExecutableSchema(graphql.Config{Resolvers: gqlResolver}),
	)

	// GraphQL endpoint (protected)
	r.POST("/graphql", middleware.AuthMiddleware(), func(c *gin.Context) {
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
