package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/status_page/backend/api"
	"github.com/status_page/backend/db"
	"github.com/status_page/backend/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Set default port if not specified
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize database connection
	db.Connect()
	db.MigrateDB()

	// Initialize Gin router
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Update in production to your domain
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Public routes
	public := r.Group("/api")
	{
		// Auth routes
		public.POST("/auth/signup", api.Signup)
		public.POST("/auth/login", api.Login)

		// Public status page routes - no authentication required
		public.GET("/public/:orgId/services", api.GetPublicServices)
		public.GET("/public/:orgId/incidents", api.GetPublicIncidents)

		// WebSocket connection for real-time updates
		public.GET("/ws/:orgId", api.HandleWebSocket)
	}

	// Protected routes - require authentication
	protected := r.Group("/api")
	protected.Use(middleware.Auth())
	{
		// Service management
		protected.GET("/services", api.GetServices)
		protected.GET("/services/:id", api.GetService)
		protected.POST("/services", api.CreateService)
		protected.PUT("/services/:id", api.UpdateService)
		protected.DELETE("/services/:id", api.DeleteService)

		// Incident management
		protected.GET("/incidents", api.GetIncidents)
		protected.GET("/incidents/:id", api.GetIncident)
		protected.POST("/incidents", api.CreateIncident)
		protected.PUT("/incidents/:id", api.UpdateIncident)
		protected.DELETE("/incidents/:id", api.DeleteIncident)

		// Incident updates
		protected.POST("/incidents/:id/updates", api.AddIncidentUpdate)
	}

	// Start the server
	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}
