package services

import (
	"auth/db/models"
	"auth/types"
	"auth/utils"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Auth

func SignUpUser(db *gorm.DB, data types.SignUpStruct) (models.User, error) {
	var user models.User

	// Check for existing user
	if err := db.Where("email = ?", data.Email).First(&user).Error; err == nil {
		// User already exists
		return models.User{}, gorm.ErrRegistered
	} else if err != gorm.ErrRecordNotFound {
		// Some other error occurred
		return models.User{}, err
	}

	// Hash the password
	hashedPass, err := utils.HashPassword(data.Password)
	if err != nil {
		return models.User{}, err
	}

	user = models.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: hashedPass,
	}

	inst := db.Create(&user)
	if inst.Error != nil {
		return models.User{}, inst.Error
	}

	return user, inst.Error
}

func SignInUser(c *gin.Context, db *gorm.DB, data types.SignInStruct) error {
	var user models.User

	// Check for existing user
	if err := db.Where("email = ?", data.Email).First(&user).Error; err != gorm.ErrRecordNotFound && err != nil {
		// Some other error occurred
		return err
	}

	doesPasswordMatch := utils.CheckPasswordHash(data.Password, user.Password)
	if !doesPasswordMatch {
		return fmt.Errorf("Invalid email or password")
	}

	// Create or find an existing session
	err := utils.CreateOrFindSession(c, db, user.ID)
	if err != nil {
		return fmt.Errorf("Error creating/finding a session")
	}

	return nil
}

// Thought

func CreateThought(c *gin.Context, db *gorm.DB, data types.CreateThought) (models.Thought, error) {
	// Get user id from the session
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		return models.Thought{}, fmt.Errorf("Unauthorized")
	}

	userID, err := utils.GetUserIDFromSessionID(db, sessionID)
	if err != nil {
		return models.Thought{}, err
	}

	thought := models.Thought{
		ID:         uuid.New(),
		UserID:     userID,
		Thought:    data.Thought,
		Visibility: data.Visibility,
	}

	inst := db.Create(&thought)
	if inst.Error != nil {
		return models.Thought{}, fmt.Errorf("Error creating a record")
	}

	return thought, nil
}

func GetUserThoughts(db *gorm.DB, userID uint) ([]models.Thought, error) {
	var thoughts []models.Thought

	if err := db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&thoughts).Error; err != nil {
		return []models.Thought{}, err
	}

	return thoughts, nil
}

func GetPublicThoughts(db *gorm.DB) ([]models.Thought, error) {
	var thoughts []models.Thought

	if err := db.
		Where("visibility = ?", "public").
		Order("created_at DESC").
		Limit(20).
		Find(&thoughts).Error; err != nil {
		return []models.Thought{}, err
	}

	return thoughts, nil
}

func UpdateThought(db *gorm.DB, updateVal types.UpdateThought, id string, userID uint) error {
	if err := db.
		Model(&models.Thought{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("visibility", updateVal.Visibility).Error; err != nil {
		return err
	}

	return nil
}
