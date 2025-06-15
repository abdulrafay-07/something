package middleware

import (
	"auth/db/models"
	"auth/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the cookie
		cookie, err := c.Request.Cookie("session_id")
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, types.ErrorResponse{
				Success: false,
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		// Parse the session id
		sessionID, err := uuid.Parse(cookie.Value)
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, types.ErrorResponse{
				Success: false,
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		// Find session from database
		var session models.Session
		if err := db.First(&session, "id = ? AND expires_at > NOW()", sessionID).Error; err != nil {
			c.IndentedJSON(http.StatusUnauthorized, types.ErrorResponse{
				Success: false,
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		c.Set("user_id", session.UserID)
		c.Next()
	}
}
