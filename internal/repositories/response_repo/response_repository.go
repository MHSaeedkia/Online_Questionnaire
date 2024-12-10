package response_repo

import (
	"gorm.io/gorm"
	"online-questionnaire/internal/models"
)

type ResponseRepository interface {
	CreateResponse(response *models.Response) error
	GetQuestionnaireResponses(questionnaireID uint) ([]*models.Response, error)
	GetResponseByID(responseID uint, response *models.Response) error
	UpdateResponse(response *models.Response) error
	GetByUserAndQuestionnaire(userID uint, questionnaireID uint) (*models.Response, error) // اضافه کردن متد جدید برای بررسی رای دادن کاربر
}

type responseRepository struct {
	db *gorm.DB
}

func NewResponseRepository(db *gorm.DB) ResponseRepository {
	return &responseRepository{db: db}
}

func (r *responseRepository) CreateResponse(response *models.Response) error {
	return r.db.Create(response).Error
}

func (r *responseRepository) GetQuestionnaireResponses(questionnaireID uint) ([]*models.Response, error) {
	var responses []*models.Response
	err := r.db.Where("questionnaire_id = ?", questionnaireID).Find(&responses).Error
	return responses, err
}

// متد جدید برای گرفتن پاسخ کاربر از یک پرسشنامه خاص
func (r *responseRepository) GetByUserAndQuestionnaire(userID uint, questionnaireID uint) (*models.Response, error) {
	var response models.Response
	err := r.db.Where("user_id = ? AND questionnaire_id = ?", userID, questionnaireID).First(&response).Error
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *responseRepository) GetResponseByID(responseID uint, response *models.Response) error {
	return r.db.First(response, "id = ?", responseID).Error
}

func (r *responseRepository) UpdateResponse(response *models.Response) error {
	return r.db.Save(response).Error
}
