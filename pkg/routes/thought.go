package routes

import (
	"auth/db/services"
	"auth/types"
	"auth/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateThoughtHandler(db *gorm.DB, c *gin.Context) {
	var thoughtData types.CreateThought

	if err := c.BindJSON(&thoughtData); err != nil {
		return
	}

	// Validation
	if len(thoughtData.Thought) < 5 {
		c.IndentedJSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Message: "Thought should has at least 5 characters",
		})
		return
	}

	if len(thoughtData.Thought) > 255 {
		c.IndentedJSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Message: "Thought cannot be greater than 255 characters",
		})
		return
	}

	thought, err := services.CreateThought(c, db, thoughtData)
	if err != nil {
		if err.Error() == "Unauthorized" {
			c.IndentedJSON(http.StatusUnauthorized, types.ErrorResponse{
				Success: false,
				Message: err.Error(),
			})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, types.ErrorResponse{
				Success: false,
				Message: err.Error(),
			})
		}
		return
	}

	resp := types.OKResponse{
		Success: true,
		Message: "Thought created successfully",
		Data: map[string]interface{}{
			"thought":    thought.Thought,
			"visibility": thought.Visibility,
		},
	}

	c.IndentedJSON(http.StatusCreated, resp)
}

func GetUserThoughtsHandler(db *gorm.DB, c *gin.Context) {
	// Get user id from the session
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, types.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	userID, err := utils.GetUserIDFromSessionID(db, sessionID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, types.ErrorResponse{
			Success: false,
			Message: "Error fetching user from session id",
		})
		return
	}

	thoughts, err := services.GetUserThoughts(db, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.IndentedJSON(http.StatusNotFound, types.ErrorResponse{
				Success: false,
				Message: "No Record found",
			})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, types.ErrorResponse{
				Success: false,
				Message: err.Error(),
			})
		}
		return
	}

	resp := types.OKResponse{
		Success: true,
		Message: "Records found",
		Data: map[string]interface{}{
			"thoughts": thoughts,
		},
	}

	c.IndentedJSON(http.StatusOK, resp)
}

func GetPublicThoughtsHandler(db *gorm.DB, c *gin.Context) {
	thoughts, err := services.GetPublicThoughts(db)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, types.ErrorResponse{
			Success: false,
			Message: "Internal server error",
		})
		return
	}

	resp := types.OKResponse{
		Success: true,
		Message: "Records found",
		Data: map[string]interface{}{
			"thoughts": thoughts,
		},
	}

	c.IndentedJSON(http.StatusOK, resp)
}

func UpdateThoughtHandler(db *gorm.DB, c *gin.Context) {
	// Get user id from the session
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, types.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	userID, err := utils.GetUserIDFromSessionID(db, sessionID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, types.ErrorResponse{
			Success: false,
			Message: "Error fetching user from session id",
		})
		return
	}

	var updateVal types.UpdateThought
	if err := c.BindJSON(&updateVal); err != nil {
		return
	}

	id := c.Param("id")

	err = services.UpdateThought(db, updateVal, id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.IndentedJSON(http.StatusNotFound, types.ErrorResponse{
				Success: false,
				Message: "Record not found",
			})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, types.ErrorResponse{
				Success: false,
				Message: err.Error(),
			})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, types.OKResponse{
		Success: true,
		Message: "Thought updated successfully",
	})
}
