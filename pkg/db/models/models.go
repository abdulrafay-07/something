package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Visibility string

const (
	VisibilityPublic  Visibility = "public"
	VisibilityPrivate Visibility = "private"
)

type User struct {
	gorm.Model
	Name     string `gorm:"size:255;not null"`
	Email    string `gorm:"size:255;not null;uniqueIndex"`
	Password string `gorm:"size:255;not null"`
}

type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uint      `gorm:"primarykey"`
	CreatedAt time.Time
	ExpiresAt time.Time
}

type Thought struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID     uint       `gorm:"primarykey"`
	Thought    string     `gorm:"type:TEXT"`
	Visibility Visibility `gorm:"type:VARCHAR(10);not null"`
	CreatedAt  time.Time
}
