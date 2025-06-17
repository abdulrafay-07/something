package utils

import (
	"auth/db/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateOrFindSession(c *gin.Context, db *gorm.DB, userId uint) error {
	var existingSession models.Session

	err := db.Where("user_id = ? AND expires_at > NOW()", userId).First(&existingSession).Error
	if err == nil {
		// Session already exists
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "session_id",
			Value:    existingSession.ID.String(),
			Path:     "/",
			HttpOnly: true,
			Secure:   false, // True for HTTPS
			SameSite: http.SameSiteLaxMode,
			Expires:  existingSession.ExpiresAt,
		})

		return nil
	} else if err == gorm.ErrRecordNotFound {
		// Session not found, create one

		newSession := models.Session{
			ID:        uuid.New(),
			UserID:    userId,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		}

		db.Create(&newSession)

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "session_id",
			Value:    newSession.ID.String(),
			Path:     "/",
			HttpOnly: true,
			Secure:   false, // True for HTTPS
			SameSite: http.SameSiteLaxMode,
			Expires:  newSession.ExpiresAt,
		})

		return nil
	}

	return err
}

func GetUserIDFromSessionID(db *gorm.DB, sessionID string) (uint, error) {
	var session models.Session

	if err := db.First(&session, "id = ? AND expires_at > NOW()", sessionID).Error; err != nil {
		return 0, err
	}

	return session.UserID, nil
}
