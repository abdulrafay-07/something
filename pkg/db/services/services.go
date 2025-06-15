package services

import (
	"auth/db/models"
	"auth/types"
	"auth/utils"
	"fmt"

	"github.com/gin-gonic/gin"
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
