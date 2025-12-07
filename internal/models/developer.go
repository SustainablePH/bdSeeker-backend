package models

import (
	"time"

	"gorm.io/gorm"
)

// DeveloperProfile represents a developer's profile
type DeveloperProfile struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;uniqueIndex" json:"user_id"`
	Bio       string         `gorm:"type:text" json:"bio"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User                 User                    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Experiences          []DeveloperExperience   `gorm:"foreignKey:DeveloperID" json:"experiences,omitempty"`
	Educations           []DeveloperEducation    `gorm:"foreignKey:DeveloperID" json:"educations,omitempty"`
	Certificates         []DeveloperCertificate  `gorm:"foreignKey:DeveloperID" json:"certificates,omitempty"`
	Technologies         []Technology            `gorm:"many2many:developer_technologies;" json:"technologies,omitempty"`
	ProgrammingLanguages []ProgrammingLanguage   `gorm:"many2many:developer_languages;" json:"programming_languages,omitempty"`
}

// DeveloperExperience represents work experience
type DeveloperExperience struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	DeveloperID uint           `gorm:"not null;index" json:"developer_id"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	CompanyName string         `gorm:"size:255;not null" json:"company_name"`
	StartDate   time.Time      `gorm:"not null" json:"start_date"`
	EndDate     *time.Time     `json:"end_date"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Developer DeveloperProfile `gorm:"foreignKey:DeveloperID" json:"developer,omitempty"`
}

// DeveloperEducation represents educational background
type DeveloperEducation struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	DeveloperID   uint           `gorm:"not null;index" json:"developer_id"`
	Institution   string         `gorm:"size:255;not null" json:"institution"`
	Degree        string         `gorm:"size:255;not null" json:"degree"`
	FieldOfStudy  string         `gorm:"size:255" json:"field_of_study"`
	StartDate     time.Time      `gorm:"not null" json:"start_date"`
	EndDate       *time.Time     `json:"end_date"`
	Grade         string         `gorm:"size:50" json:"grade"`
	Description   string         `gorm:"type:text" json:"description"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Developer DeveloperProfile `gorm:"foreignKey:DeveloperID" json:"developer,omitempty"`
}

// DeveloperCertificate represents professional certificates
type DeveloperCertificate struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	DeveloperID         uint           `gorm:"not null;index" json:"developer_id"`
	CertificateName     string         `gorm:"size:255;not null" json:"certificate_name"`
	IssuingOrganization string         `gorm:"size:255;not null" json:"issuing_organization"`
	IssueDate           time.Time      `gorm:"not null" json:"issue_date"`
	ExpirationDate      *time.Time     `json:"expiration_date"`
	CredentialID        string         `gorm:"size:255" json:"credential_id"`
	CertificateLink     string         `gorm:"size:500" json:"certificate_link"`
	Description         string         `gorm:"type:text" json:"description"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Developer DeveloperProfile `gorm:"foreignKey:DeveloperID" json:"developer,omitempty"`
}
