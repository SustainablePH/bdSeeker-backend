package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents the main user entity
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Email        string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	FullName     string         `gorm:"size:255" json:"full_name"`
	Role         string         `gorm:"size:50;not null;default:'developer'" json:"role"` // developer, company, admin
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	CompanyProfile   *CompanyProfile   `gorm:"foreignKey:UserID" json:"company_profile,omitempty"`
	DeveloperProfile *DeveloperProfile `gorm:"foreignKey:UserID" json:"developer_profile,omitempty"`
}

// Technology represents a technology/skill
type Technology struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null;size:100" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProgrammingLanguage represents a programming language
type ProgrammingLanguage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null;size:100" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
