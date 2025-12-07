package repositories

import (
	"github.com/bishworup11/bdSeeker-backend/internal/models"
	"gorm.io/gorm"
)

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) Create(company *models.CompanyProfile) error {
	return r.db.Create(company).Error
}

func (r *CompanyRepository) FindByID(id uint) (*models.CompanyProfile, error) {
	var company models.CompanyProfile
	err := r.db.Preload("User").Preload("Technologies").Preload("Ratings").
		Preload("Reviews", "is_approved = ?", true).First(&company, id).Error
	return &company, err
}

func (r *CompanyRepository) FindByUserID(userID uint) (*models.CompanyProfile, error) {
	var company models.CompanyProfile
	err := r.db.Where("user_id = ?", userID).Preload("Technologies").First(&company).Error
	return &company, err
}

func (r *CompanyRepository) Update(company *models.CompanyProfile) error {
	return r.db.Save(company).Error
}

func (r *CompanyRepository) Delete(id uint) error {
	return r.db.Delete(&models.CompanyProfile{}, id).Error
}

func (r *CompanyRepository) List(page, limit int, location string, techIDs []uint) ([]models.CompanyProfile, int64, error) {
	var companies []models.CompanyProfile
	var total int64

	offset := (page - 1) * limit
	query := r.db.Model(&models.CompanyProfile{})

	if location != "" {
		query = query.Where("location ILIKE ?", "%"+location+"%")
	}

	if len(techIDs) > 0 {
		query = query.Joins("JOIN company_technologies ON company_technologies.company_profile_id = company_profiles.id").
			Where("company_technologies.technology_id IN ?", techIDs).
			Distinct()
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Preload("Technologies").Preload("User").Find(&companies).Error
	return companies, total, err
}

// Rating operations
func (r *CompanyRepository) CreateRating(rating *models.CompanyRating) error {
	return r.db.Create(rating).Error
}

func (r *CompanyRepository) UpdateRating(rating *models.CompanyRating) error {
	return r.db.Save(rating).Error
}

func (r *CompanyRepository) FindRatingByUserAndCompany(userID, companyID uint) (*models.CompanyRating, error) {
	var rating models.CompanyRating
	err := r.db.Where("user_id = ? AND company_id = ?", userID, companyID).First(&rating).Error
	return &rating, err
}

func (r *CompanyRepository) GetAverageRating(companyID uint) (float64, error) {
	var avg float64
	err := r.db.Model(&models.CompanyRating{}).Where("company_id = ?", companyID).
		Select("COALESCE(AVG(rating), 0)").Scan(&avg).Error
	return avg, err
}

// Review operations
func (r *CompanyRepository) CreateReview(review *models.CompanyReview) error {
	return r.db.Create(review).Error
}

func (r *CompanyRepository) FindReviewByID(id uint) (*models.CompanyReview, error) {
	var review models.CompanyReview
	err := r.db.Preload("User").Preload("Reactions").Preload("Comments").First(&review, id).Error
	return &review, err
}

func (r *CompanyRepository) UpdateReview(review *models.CompanyReview) error {
	return r.db.Save(review).Error
}

func (r *CompanyRepository) DeleteReview(id uint) error {
	return r.db.Delete(&models.CompanyReview{}, id).Error
}

func (r *CompanyRepository) ListReviews(companyID uint, page, limit int, approvedOnly bool) ([]models.CompanyReview, int64, error) {
	var reviews []models.CompanyReview
	var total int64

	offset := (page - 1) * limit
	query := r.db.Model(&models.CompanyReview{}).Where("company_id = ?", companyID)

	if approvedOnly {
		query = query.Where("is_approved = ?", true)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Preload("User").Preload("Reactions").
		Preload("Comments", "is_approved = ?", true).Find(&reviews).Error
	return reviews, total, err
}

// Review reaction operations
func (r *CompanyRepository) CreateReviewReaction(reaction *models.CompanyReviewReaction) error {
	return r.db.Create(reaction).Error
}

func (r *CompanyRepository) UpdateReviewReaction(reaction *models.CompanyReviewReaction) error {
	return r.db.Save(reaction).Error
}

func (r *CompanyRepository) FindReviewReaction(reviewID, userID uint) (*models.CompanyReviewReaction, error) {
	var reaction models.CompanyReviewReaction
	err := r.db.Where("review_id = ? AND user_id = ?", reviewID, userID).First(&reaction).Error
	return &reaction, err
}

// Review comment operations
func (r *CompanyRepository) CreateReviewComment(comment *models.CompanyReviewComment) error {
	return r.db.Create(comment).Error
}

func (r *CompanyRepository) FindReviewCommentByID(id uint) (*models.CompanyReviewComment, error) {
	var comment models.CompanyReviewComment
	err := r.db.Preload("Replies").First(&comment, id).Error
	return &comment, err
}

func (r *CompanyRepository) UpdateReviewComment(comment *models.CompanyReviewComment) error {
	return r.db.Save(comment).Error
}

// Review reply operations
func (r *CompanyRepository) CreateReviewReply(reply *models.CompanyReviewReply) error {
	return r.db.Create(reply).Error
}
