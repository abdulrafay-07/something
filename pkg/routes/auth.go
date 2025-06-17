package routes

import (
	"auth/db/models"
	"auth/db/services"
	"auth/types"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SignUpHandler(db *gorm.DB, c *gin.Context) {
	var signUpData types.SignUpStruct

	if err := c.BindJSON(&signUpData); err != nil {
		return
	}

	// Validation
	if signUpData.Name == "" || signUpData.Email == "" {
		c.IndentedJSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Message: "Name and Email are required",
		})
		return
	}

	if signUpData.Password == "" || len(signUpData.Password) < 6 {
		c.IndentedJSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Message: "Password must be at least 6 characters long",
		})
		return
	}

	// Hash password and store in database
	user, err := services.SignUpUser(db, signUpData)
	if err != nil {
		log.Printf("services.SignUpUser Error: %v", err)
		if err == gorm.ErrRegistered {
			c.IndentedJSON(http.StatusOK, types.ErrorResponse{
				Success: false,
				Message: "User already exists",
			})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, types.ErrorResponse{
				Success: false,
				Message: "Internal server error",
			})
		}
		return
	}

	resp := types.OKResponse{
		Success: true,
		Message: "Account created",
		Data: types.PublicUser{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}
	c.IndentedJSON(http.StatusCreated, resp)
}

func SignInHandler(db *gorm.DB, c *gin.Context) {
	var signInData types.SignInStruct

	if err := c.BindJSON(&signInData); err != nil {
		return
	}

	// Validation
	if signInData.Email == "" {
		c.IndentedJSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Message: "Email field is required",
		})
		return
	}

	if signInData.Password == "" || len(signInData.Password) < 6 {
		c.IndentedJSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Message: "Password must be at least 6 characters long",
		})
		return
	}

	// Database logic
	err := services.SignInUser(c, db, signInData)
	if err != nil {
		log.Printf("services.SignInUser Error: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, types.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	resp := types.OKResponse{
		Success: true,
		Message: "Logged in successfully",
	}

	c.IndentedJSON(http.StatusOK, resp)
}

func LogoutHandler(db *gorm.DB, c *gin.Context) {
	sessionIDCookie, _ := c.Request.Cookie("session_id")
	sessionID := sessionIDCookie.Value

	// Delete the session from the database
	db.Delete(&models.Session{}, "id = ?", sessionID)

	// Clear the cookie from the header
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	})

	c.IndentedJSON(http.StatusOK, types.OKResponse{
		Success: true,
		Message: "Logged out",
	})
}

func MeHandler(db *gorm.DB, c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.IndentedJSON(http.StatusOK, types.ErrorResponse{
			Success: false,
			Message: "Cookie not found",
		})
		return
	}

	// Find session from database
	var session models.Session
	if err := db.First(&session, "id = ? AND expires_at > NOW()", sessionID).Error; err != nil {
		c.IndentedJSON(http.StatusOK, types.ErrorResponse{
			Success: false,
			Message: "Session not found",
		})
		return
	}

	// Find the user
	var user models.User
	if err := db.First(&user, "id = ?", session.UserID).Error; err != nil {
		c.IndentedJSON(http.StatusOK, types.ErrorResponse{
			Success: false,
			Message: "User not found",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, types.OKResponse{
		Success: true,
		Message: "Authorized",
		Data: types.PublicUser{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}
