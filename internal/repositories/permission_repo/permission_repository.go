package permission_repo

import (
	"online-questionnaire/internal/models"
	"time"

	"gorm.io/gorm"
)

type PermissionRepository interface {
	GrantPermission(questionnaireID uint, userID uint, permission models.Name, expiresAt *string) error
	GetPendingPermissions(questionnaireID uint) ([]models.QuestionnairePermission, error)
	GetQuestionnairePermission(id uint) (models.QuestionnairePermission, error)
	UpdateQuestionnairePermission(permission *models.QuestionnairePermission) error
	CreateQuestionnairePermission(permission *models.QuestionnairePermission) error
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db}
}

func (r *permissionRepository) GrantPermission(questionnaireID uint, userID uint, permission models.Name, expiresAt *string) error {
	permissionRecord := models.QuestionnairePermission{
		QuestionnaireID: questionnaireID,
		UserID:          userID,
		Permission:      models.Permission{Name: permission},
	}

	if expiresAt != nil {
		expTime, err := time.Parse(time.RFC3339, *expiresAt)
		if err != nil {
			return err
		}
		permissionRecord.ExpiresAt = &expTime
	}

	return r.db.Create(&permissionRecord).Error
}

func (r *permissionRepository) GetPendingPermissions(questionnaireID uint) ([]models.QuestionnairePermission, error) {
	var requests []models.QuestionnairePermission
	err := r.db.Where("questionnaire_id = ? AND status = ?", questionnaireID, "Pending").Find(&requests).Error
	return requests, err
}

func (r *permissionRepository) GetQuestionnairePermission(id uint) (models.QuestionnairePermission, error) {
	var permission models.QuestionnairePermission
	err := r.db.First(&permission, id).Error
	return permission, err
}

func (r *permissionRepository) UpdateQuestionnairePermission(permission *models.QuestionnairePermission) error {
	return r.db.Save(permission).Error
}

func (r *permissionRepository) CreateQuestionnairePermission(permission *models.QuestionnairePermission) error {
	return r.db.Create(permission).Error
}
