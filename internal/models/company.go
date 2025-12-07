package models

import (
	"time"

	"gorm.io/gorm"
)

// CompanyProfile represents a company's profile
type CompanyProfile struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;uniqueIndex" json:"user_id"`
	CompanyName string         `gorm:"size:255;not null" json:"company_name"`
	Description string         `gorm:"type:text" json:"description"`
	Website     string         `gorm:"size:255" json:"website"`
	Location    string         `gorm:"size:255" json:"location"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User         User                    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Technologies []Technology            `gorm:"many2many:company_technologies;" json:"technologies,omitempty"`
	JobPosts     []JobPost               `gorm:"foreignKey:CompanyID" json:"job_posts,omitempty"`
	Ratings      []CompanyRating         `gorm:"foreignKey:CompanyID" json:"ratings,omitempty"`
	Reviews      []CompanyReview         `gorm:"foreignKey:CompanyID" json:"reviews,omitempty"`
}

// CompanyRating represents a rating given to a company
type CompanyRating struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CompanyID uint      `gorm:"not null;index" json:"company_id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Rating    int       `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Company CompanyProfile `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	User    User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// CompanyReview represents a review for a company
type CompanyReview struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CompanyID  uint           `gorm:"not null;index" json:"company_id"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	IsApproved bool           `gorm:"default:false" json:"is_approved"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Company   CompanyProfile            `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	User      User                      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Reactions []CompanyReviewReaction   `gorm:"foreignKey:ReviewID" json:"reactions,omitempty"`
	Comments  []CompanyReviewComment    `gorm:"foreignKey:ReviewID" json:"comments,omitempty"`
}

// CompanyReviewReaction represents a reaction to a company review
type CompanyReviewReaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ReviewID  uint      `gorm:"not null;index" json:"review_id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Type      string    `gorm:"size:50;not null" json:"type"` // like, useful
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Review CompanyReview `gorm:"foreignKey:ReviewID" json:"review,omitempty"`
	User   User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// CompanyReviewComment represents a comment on a company review
type CompanyReviewComment struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ReviewID   uint           `gorm:"not null;index" json:"review_id"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	IsApproved bool           `gorm:"default:false" json:"is_approved"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Review  CompanyReview           `gorm:"foreignKey:ReviewID" json:"review,omitempty"`
	User    User                    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Replies []CompanyReviewReply    `gorm:"foreignKey:CommentID" json:"replies,omitempty"`
}

// CompanyReviewReply represents a reply to a comment
type CompanyReviewReply struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CommentID uint      `gorm:"not null;index" json:"comment_id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Comment CompanyReviewComment `gorm:"foreignKey:CommentID" json:"comment,omitempty"`
	User    User                 `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
