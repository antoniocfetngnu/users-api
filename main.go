package main

import (
	"log"
	"net/http"
	"github.com/antoniocfetngnu/users-api/database"
	"github.com/antoniocfetngnu/users-api/graphql"

	"github.com/gin-gonic/gin"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Create GraphQL resolver (no DB field needed)
	resolver := &graphql.Resolver{}
	
	// Create GraphQL handler
	gqlHandler := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))
	
	// Setup Gin router
	r := gin.Default()

	// REST Health endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"message": "Users API is running",
			"version": "1.0.0",
		})
	})

	// GraphQL endpoint
	r.POST("/graphql", func(c *gin.Context) {
		gqlHandler.ServeHTTP(c.Writer, c.Request)
	})

	// GraphQL Playground (UI for testing)
	r.GET("/playground", func(c *gin.Context) {
		playground.Handler("GraphQL Playground", "/graphql").ServeHTTP(c.Writer, c.Request)
	})

	// Basic REST users endpoint (preview)
	r.GET("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "REST users endpoint - use GraphQL for full functionality",
		})
	})

	log.Println("🚀 Server running on http://localhost:8080")
	log.Println("📊 Health check: http://localhost:8080/health")
	log.Println("🎮 GraphQL Playground: http://localhost:8080/playground")
	
	r.Run(":8080")
}