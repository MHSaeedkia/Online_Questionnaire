package repositories

import (
	"gorm.io/gorm"
	"online-questionnaire/internal/models"
)

type ResponseRepository interface {
	CreateResponse(response *models.Response) error
	GetQuestionnaireResponses(questionnaireID uint) ([]*models.Response, error)
	GetResponseByID(responseID uint, response *models.Response) error
	UpdateResponse(response *models.Response) error
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

// UpdateResponse updates an existing response.
func (r *responseRepository) UpdateResponse(response *models.Response) error {
	return r.db.Save(response).Error
}
