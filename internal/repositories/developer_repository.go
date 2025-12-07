package repositories

import (
	"github.com/bishworup11/bdSeeker-backend/internal/models"
	"gorm.io/gorm"
)

type DeveloperRepository struct {
	db *gorm.DB
}

func NewDeveloperRepository(db *gorm.DB) *DeveloperRepository {
	return &DeveloperRepository{db: db}
}

func (r *DeveloperRepository) Create(developer *models.DeveloperProfile) error {
	return r.db.Create(developer).Error
}

func (r *DeveloperRepository) FindByID(id uint) (*models.DeveloperProfile, error) {
	var developer models.DeveloperProfile
	err := r.db.Preload("User").Preload("Technologies").Preload("ProgrammingLanguages").
		Preload("Experiences").Preload("Educations").Preload("Certificates").First(&developer, id).Error
	return &developer, err
}

func (r *DeveloperRepository) FindByUserID(userID uint) (*models.DeveloperProfile, error) {
	var developer models.DeveloperProfile
	err := r.db.Where("user_id = ?", userID).Preload("Technologies").Preload("ProgrammingLanguages").First(&developer).Error
	return &developer, err
}

func (r *DeveloperRepository) Update(developer *models.DeveloperProfile) error {
	return r.db.Save(developer).Error
}

func (r *DeveloperRepository) Delete(id uint) error {
	return r.db.Delete(&models.DeveloperProfile{}, id).Error
}

func (r *DeveloperRepository) List(page, limit int, techIDs, langIDs []uint) ([]models.DeveloperProfile, int64, error) {
	var developers []models.DeveloperProfile
	var total int64

	offset := (page - 1) * limit
	query := r.db.Model(&models.DeveloperProfile{})

	if len(techIDs) > 0 {
		query = query.Joins("JOIN developer_technologies ON developer_technologies.developer_profile_id = developer_profiles.id").
			Where("developer_technologies.technology_id IN ?", techIDs).
			Distinct()
	}

	if len(langIDs) > 0 {
		query = query.Joins("JOIN developer_languages ON developer_languages.developer_profile_id = developer_profiles.id").
			Where("developer_languages.programming_language_id IN ?", langIDs).
			Distinct()
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Preload("User").Preload("Technologies").
		Preload("ProgrammingLanguages").Find(&developers).Error
	return developers, total, err
}

// Experience operations
func (r *DeveloperRepository) CreateExperience(exp *models.DeveloperExperience) error {
	return r.db.Create(exp).Error
}

func (r *DeveloperRepository) FindExperienceByID(id uint) (*models.DeveloperExperience, error) {
	var exp models.DeveloperExperience
	err := r.db.First(&exp, id).Error
	return &exp, err
}

func (r *DeveloperRepository) UpdateExperience(exp *models.DeveloperExperience) error {
	return r.db.Save(exp).Error
}

func (r *DeveloperRepository) DeleteExperience(id uint) error {
	return r.db.Delete(&models.DeveloperExperience{}, id).Error
}

// Education operations
func (r *DeveloperRepository) CreateEducation(edu *models.DeveloperEducation) error {
	return r.db.Create(edu).Error
}

func (r *DeveloperRepository) FindEducationByID(id uint) (*models.DeveloperEducation, error) {
	var edu models.DeveloperEducation
	err := r.db.First(&edu, id).Error
	return &edu, err
}

func (r *DeveloperRepository) UpdateEducation(edu *models.DeveloperEducation) error {
	return r.db.Save(edu).Error
}

func (r *DeveloperRepository) DeleteEducation(id uint) error {
	return r.db.Delete(&models.DeveloperEducation{}, id).Error
}

// Certificate operations
func (r *DeveloperRepository) CreateCertificate(cert *models.DeveloperCertificate) error {
	return r.db.Create(cert).Error
}

func (r *DeveloperRepository) FindCertificateByID(id uint) (*models.DeveloperCertificate, error) {
	var cert models.DeveloperCertificate
	err := r.db.First(&cert, id).Error
	return &cert, err
}

func (r *DeveloperRepository) UpdateCertificate(cert *models.DeveloperCertificate) error {
	return r.db.Save(cert).Error
}

func (r *DeveloperRepository) DeleteCertificate(id uint) error {
	return r.db.Delete(&models.DeveloperCertificate{}, id).Error
}
