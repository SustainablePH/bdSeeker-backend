package repositories

import (
	"github.com/bishworup11/bdSeeker-backend/internal/models"
	"gorm.io/gorm"
)

type JobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) Create(job *models.JobPost) error {
	return r.db.Create(job).Error
}

func (r *JobRepository) FindByID(id uint) (*models.JobPost, error) {
	var job models.JobPost
	err := r.db.Preload("Company.User").Preload("Reactions").
		Preload("Comments.User").Preload("Comments.Replies.User").First(&job, id).Error
	return &job, err
}

func (r *JobRepository) Update(job *models.JobPost) error {
	return r.db.Save(job).Error
}

func (r *JobRepository) Delete(id uint) error {
	return r.db.Delete(&models.JobPost{}, id).Error
}

func (r *JobRepository) List(page, limit int, filters map[string]interface{}) ([]models.JobPost, int64, error) {
	var jobs []models.JobPost
	var total int64

	offset := (page - 1) * limit
	query := r.db.Model(&models.JobPost{})

	// Apply filters
	if companyID, ok := filters["company_id"].(uint); ok && companyID > 0 {
		query = query.Where("company_id = ?", companyID)
	}

	if workMode, ok := filters["work_mode"].(string); ok && workMode != "" {
		query = query.Where("work_mode = ?", workMode)
	}

	if location, ok := filters["location"].(string); ok && location != "" {
		query = query.Where("location ILIKE ?", "%"+location+"%")
	}

	if search, ok := filters["search"].(string); ok && search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if minSalary, ok := filters["min_salary"].(float64); ok && minSalary > 0 {
		query = query.Where("salary_max >= ?", minSalary)
	}

	if maxSalary, ok := filters["max_salary"].(float64); ok && maxSalary > 0 {
		query = query.Where("salary_min <= ?", maxSalary)
	}

	if minExp, ok := filters["min_experience"].(int); ok && minExp >= 0 {
		query = query.Where("experience_max_years >= ?", minExp)
	}

	if maxExp, ok := filters["max_experience"].(int); ok && maxExp > 0 {
		query = query.Where("experience_min_years <= ?", maxExp)
	}

	// Sorting
	if sortBy, ok := filters["sort_by"].(string); ok && sortBy != "" {
		switch sortBy {
		case "salary_desc":
			query = query.Order("salary_max DESC")
		case "salary_asc":
			query = query.Order("salary_min ASC")
		case "created_desc":
			query = query.Order("created_at DESC")
		case "created_asc":
			query = query.Order("created_at ASC")
		default:
			query = query.Order("created_at DESC")
		}
	} else {
		query = query.Order("created_at DESC")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Preload("Company.User").Find(&jobs).Error
	return jobs, total, err
}

// Reaction operations
func (r *JobRepository) CreateReaction(reaction *models.PostReaction) error {
	return r.db.Create(reaction).Error
}

func (r *JobRepository) UpdateReaction(reaction *models.PostReaction) error {
	return r.db.Save(reaction).Error
}

func (r *JobRepository) FindReaction(userID, jobPostID uint, reactionType string) (*models.PostReaction, error) {
	var reaction models.PostReaction
	query := r.db.Where("user_id = ? AND job_post_id = ?", userID, jobPostID)
	if reactionType != "" {
		query = query.Where("type = ?", reactionType)
	}
	err := query.First(&reaction).Error
	return &reaction, err
}

func (r *JobRepository) DeleteReaction(userID, jobPostID uint) error {
	return r.db.Where("user_id = ? AND job_post_id = ?", userID, jobPostID).Delete(&models.PostReaction{}).Error
}

// Comment operations
func (r *JobRepository) CreateComment(comment *models.PostComment) error {
	return r.db.Create(comment).Error
}

func (r *JobRepository) FindCommentByID(id uint) (*models.PostComment, error) {
	var comment models.PostComment
	err := r.db.Preload("User").Preload("Replies.User").First(&comment, id).Error
	return &comment, err
}

func (r *JobRepository) UpdateComment(comment *models.PostComment) error {
	return r.db.Save(comment).Error
}

func (r *JobRepository) DeleteComment(id uint) error {
	return r.db.Delete(&models.PostComment{}, id).Error
}

// Reply operations
func (r *JobRepository) CreateReply(reply *models.CommentReply) error {
	return r.db.Create(reply).Error
}
