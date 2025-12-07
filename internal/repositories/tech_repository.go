package repositories

import (
	"github.com/bishworup11/bdSeeker-backend/internal/models"
	"gorm.io/gorm"
)

type TechRepository struct {
	db *gorm.DB
}

func NewTechRepository(db *gorm.DB) *TechRepository {
	return &TechRepository{db: db}
}

// Technology operations
func (r *TechRepository) CreateTechnology(tech *models.Technology) error {
	return r.db.Create(tech).Error
}

func (r *TechRepository) FindTechnologyByID(id uint) (*models.Technology, error) {
	var tech models.Technology
	err := r.db.First(&tech, id).Error
	return &tech, err
}

func (r *TechRepository) FindTechnologyByName(name string) (*models.Technology, error) {
	var tech models.Technology
	err := r.db.Where("name = ?", name).First(&tech).Error
	return &tech, err
}

func (r *TechRepository) ListTechnologies(search string) ([]models.Technology, error) {
	var techs []models.Technology
	query := r.db.Model(&models.Technology{})

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	err := query.Order("name ASC").Find(&techs).Error
	return techs, err
}

// Programming Language operations
func (r *TechRepository) CreateLanguage(lang *models.ProgrammingLanguage) error {
	return r.db.Create(lang).Error
}

func (r *TechRepository) FindLanguageByID(id uint) (*models.ProgrammingLanguage, error) {
	var lang models.ProgrammingLanguage
	err := r.db.First(&lang, id).Error
	return &lang, err
}

func (r *TechRepository) FindLanguageByName(name string) (*models.ProgrammingLanguage, error) {
	var lang models.ProgrammingLanguage
	err := r.db.Where("name = ?", name).First(&lang).Error
	return &lang, err
}

func (r *TechRepository) ListLanguages(search string) ([]models.ProgrammingLanguage, error) {
	var langs []models.ProgrammingLanguage
	query := r.db.Model(&models.ProgrammingLanguage{})

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	err := query.Order("name ASC").Find(&langs).Error
	return langs, err
}
