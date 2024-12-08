package repositories

import (
	"online-questionnaire/internal/models"

	"gorm.io/gorm"
)

type ResponseRepository interface {
	GetQuestionnaireResponses(questionnaireID uint) ([]*models.Response, error)
	GetResponseByID(responseID uint, response *models.Response) error
	CreateResponse(response *models.Response) error
	GetResponse(responseId uint) (*models.Response, error)
	UpdateResponse(response *models.Response) error
	DeleteResponse(responseId uint) error
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

// GetResponseByID fetches a response by its ID.
func (r *responseRepository) GetResponseByID(responseID uint, response *models.Response) error {
	return r.db.First(response, "id = ?", responseID).Error
}

func (r *responseRepository) GetResponse(responseID uint) (*models.Response, error) {
	var answer *models.Response
	return answer, r.db.Find(&answer).Where("id = ?", responseID).Error
}

func (r *responseRepository) UpdateResponse(response *models.Response) error {
	return r.db.Updates(response).Where("id = ?", response.ID).Error
}

func (r *responseRepository) DeleteResponse(responseID uint) error {
	return r.db.Delete(&models.Response{}).Where("id = ?", responseID).Error
}
