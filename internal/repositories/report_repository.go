package repositories

import (
	"github.com/bishworup11/bdSeeker-backend/internal/models"
	"gorm.io/gorm"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) Create(report *models.UserReport) error {
	return r.db.Create(report).Error
}

func (r *ReportRepository) FindByID(id uint) (*models.UserReport, error) {
	var report models.UserReport
	err := r.db.Preload("Reporter").Preload("Reported").Preload("Reviewer").First(&report, id).Error
	return &report, err
}

func (r *ReportRepository) Update(report *models.UserReport) error {
	return r.db.Save(report).Error
}

func (r *ReportRepository) List(page, limit int, status, reportType string) ([]models.UserReport, int64, error) {
	var reports []models.UserReport
	var total int64

	offset := (page - 1) * limit
	query := r.db.Model(&models.UserReport{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if reportType != "" {
		query = query.Where("report_type = ?", reportType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Preload("Reporter").Preload("Reported").
		Preload("Reviewer").Order("created_at DESC").Find(&reports).Error
	return reports, total, err
}
