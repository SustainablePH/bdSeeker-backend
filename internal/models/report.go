package models

import (
	"time"

	"gorm.io/gorm"
)

// UserReport represents a report filed by one user against another
type UserReport struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ReporterID uint           `gorm:"not null;index" json:"reporter_id"`
	ReportedID uint           `gorm:"not null;index" json:"reported_id"`
	ReportType string         `gorm:"size:100;not null" json:"report_type"` // harassment, spam, inappropriate, etc.
	Description string        `gorm:"type:text;not null" json:"description"`
	Status     string         `gorm:"size:50;not null;default:'pending'" json:"status"` // pending, reviewed, resolved, dismissed
	ReviewedBy *uint          `gorm:"index" json:"reviewed_by"`
	ReviewedAt *time.Time     `json:"reviewed_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Reporter User  `gorm:"foreignKey:ReporterID" json:"reporter,omitempty"`
	Reported User  `gorm:"foreignKey:ReportedID" json:"reported,omitempty"`
	Reviewer *User `gorm:"foreignKey:ReviewedBy" json:"reviewer,omitempty"`
}
