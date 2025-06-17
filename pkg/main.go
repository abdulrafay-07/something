package main

import (
	"auth/db"
	"auth/db/models"
	"auth/middleware"
	"auth/routes"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api")
	{
		api.POST("/signup", func(c *gin.Context) {
			routes.SignUpHandler(db, c)
		})
		api.POST("/signin", func(c *gin.Context) {
			routes.SignInHandler(db, c)
		})
		api.GET("/me", func(c *gin.Context) {
			routes.MeHandler(db, c)
		})
		api.POST("/logout", middleware.AuthMiddleware(db), func(c *gin.Context) {
			routes.LogoutHandler(db, c)
		})

		api.POST("/thought", middleware.AuthMiddleware(db), func(c *gin.Context) {
			routes.CreateThoughtHandler(db, c)
		})
		api.GET("/thought", middleware.AuthMiddleware(db), func(c *gin.Context) {
			routes.GetUserThoughtsHandler(db, c)
		})
		api.GET("/thoughts/public", middleware.AuthMiddleware(db), func(c *gin.Context) {
			routes.GetPublicThoughtsHandler(db, c)
		})
		api.PATCH("/thought/:id", middleware.AuthMiddleware(db), func(c *gin.Context) {
			routes.UpdateThoughtHandler(db, c)
		})
	}
}

func main() {
	database := db.ConnectDB()

	// Auto migrate the schemas
	err := database.AutoMigrate(&models.User{}, &models.Session{}, &models.Thought{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Migrated DB schema")

	// Setup router
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	SetRoutes(router, database)

	router.Run("localhost:8080")
}
