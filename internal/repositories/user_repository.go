package repositories

import (
	"online-questionnaire/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Quesionnare(ownerId uint) ([]models.Questionnaire, error)
	EditQuestionnare(ownerId uint, questionnareId uint, questionnare *models.Questionnaire) error
	CancleQuestionnare(questionnareId uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Quesionnare(ownerId uint) ([]models.Questionnaire, error) {
	var questionnares []models.Questionnaire
	return questionnares, r.db.Where("owner_id = ?", ownerId).Find(&questionnares).Error
}

func (r *userRepository) EditQuestionnare(ownerId uint, questionnareId uint, questionnare *models.Questionnaire) error {
	return r.db.Where("owner_id = ? AND id = ?", ownerId, questionnareId).Updates(&questionnare).Error
}

func (r *userRepository) CancleQuestionnare(quuestionnareId uint) error {
	return r.db.Where("id = ?", quuestionnareId).Delete(models.Questionnaire{}).Error
}
