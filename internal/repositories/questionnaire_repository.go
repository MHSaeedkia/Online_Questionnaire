package repositories

import (
	"online-questionnaire/internal/models"

	"gorm.io/gorm"
)

type QuestionnaireRepository interface {
	CreateQuestionnaire(questionnaire *models.Questionnaire) error
	CreateQuestion(quesion *models.Question) error
	GetQuestion(quesionId uint) (*models.Question, error)
	UpdateQuestion(quesion *models.Question) error
	DeleteQuestion(quesionId uint) error
	CreateAnswer(answer *models.Response) error
	GetAnswer(answerId uint) (*models.Response, error)
	UpdateAnswer(answer *models.Response) error
	DeleteAnswer(answerId uint) error
}

type questionnaireRepository struct {
	db *gorm.DB
}

func NewQuestionnaireRepository(db *gorm.DB) QuestionnaireRepository {
	return &questionnaireRepository{db}
}

func (r *questionnaireRepository) CreateQuestionnaire(questionnaire *models.Questionnaire) error {
	return r.db.Create(questionnaire).Error
}

func (r *questionnaireRepository) CreateQuestion(quesion *models.Question) error {
	return r.db.Create(quesion).Error
}

func (r *questionnaireRepository) GetQuestion(quesionId uint) (*models.Question, error) {
	var question *models.Question
	return question, r.db.Find(&question).Where("id = ?", quesionId).Error
}

func (r *questionnaireRepository) UpdateQuestion(quesion *models.Question) error {
	return r.db.Updates(quesion).Where("id = ?", quesion.ID).Error
}

func (r *questionnaireRepository) DeleteQuestion(quesionId uint) error {
	return r.db.Delete(&models.Question{}).Where("id = ?", quesionId).Error
}

func (r *questionnaireRepository) CreateAnswer(response *models.Response) error {
	return r.db.Create(response).Error
}

func (r *questionnaireRepository) GetAnswer(answerId uint) (*models.Response, error) {
	var answer *models.Response
	return answer, r.db.Find(&answer).Where("id = ?", answerId).Error
}

func (r *questionnaireRepository) UpdateAnswer(answer *models.Response) error {
	return r.db.Updates(answer).Where("id = ?", answer.ID).Error
}

func (r *questionnaireRepository) DeleteAnswer(answerId uint) error {
	return r.db.Delete(&models.Response{}).Where("id = ?", answerId).Error
}
