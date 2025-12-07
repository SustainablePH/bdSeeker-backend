package models

import (
	"time"

	"gorm.io/gorm"
)

// JobPost represents a job posting
type JobPost struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	CompanyID         uint           `gorm:"not null;index" json:"company_id"`
	Title             string         `gorm:"size:255;not null" json:"title"`
	Description       string         `gorm:"type:text;not null" json:"description"`
	SalaryMin         float64        `json:"salary_min"`
	SalaryMax         float64        `json:"salary_max"`
	ExperienceMinYears int           `json:"experience_min_years"`
	ExperienceMaxYears int           `json:"experience_max_years"`
	WorkMode          string         `gorm:"size:50" json:"work_mode"` // office, hybrid, remote, onsite
	Location          string         `gorm:"size:255" json:"location"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Company   CompanyProfile `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Reactions []PostReaction `gorm:"foreignKey:JobPostID" json:"reactions,omitempty"`
	Comments  []PostComment  `gorm:"foreignKey:JobPostID" json:"comments,omitempty"`
}

// PostReaction represents a reaction to a job post
type PostReaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	JobPostID uint      `gorm:"not null;index" json:"job_post_id"`
	Type      string    `gorm:"size:50;not null" json:"type"` // like, bookmark, apply
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	JobPost JobPost `gorm:"foreignKey:JobPostID" json:"job_post,omitempty"`
}

// PostComment represents a comment on a job post
type PostComment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	JobPostID uint           `gorm:"not null;index" json:"job_post_id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User    User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	JobPost JobPost        `gorm:"foreignKey:JobPostID" json:"job_post,omitempty"`
	Replies []CommentReply `gorm:"foreignKey:CommentID" json:"replies,omitempty"`
}

// CommentReply represents a reply to a comment
type CommentReply struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	CommentID uint        `gorm:"not null;index" json:"comment_id"`
	UserID    uint        `gorm:"not null;index" json:"user_id"`
	Content   string      `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`

	// Relations
	Comment PostComment `gorm:"foreignKey:CommentID" json:"comment,omitempty"`
	User    User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
