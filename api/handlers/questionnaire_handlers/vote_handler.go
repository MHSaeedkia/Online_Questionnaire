package questionnaire_handlers

import (
    "online-questionnaire/internal/models"
    "online-questionnaire/internal/repositories/question_repo"
    "online-questionnaire/internal/repositories/response_repo"
    "online-questionnaire/internal/repositories/questionnaire_repo"
    "github.com/gofiber/fiber/v2"
    "strconv"
    "time"
)

type VoteHandler struct {
    questionnaireRepo *questionnaire_repo.QuestionnaireRepository
    questionRepo      *question_repo.QuestionRepository
    responseRepo      *response_repo.ResponseRepository
}

func NewVoteHandler(qRepo *questionnaire_repo.QuestionnaireRepository, qnRepo *question_repo.QuestionRepository, rRepo *response_repo.ResponseRepository) *VoteHandler {
    return &VoteHandler{
        questionnaireRepo: qRepo,
        questionRepo:      qnRepo,
        responseRepo:      rRepo,
    }
}

func (h *VoteHandler) VoteOnQuestionnaire(c *fiber.Ctx) error {
    questionnaireID, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid questionnaire ID"})
    }

    questionnaire, err := h.questionnaireRepo.GetByID(uint(questionnaireID))
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Questionnaire not found"})
    }

    if time.Now().After(questionnaire.EndTime) {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Voting period has ended"})
    }

    var userResponse models.Response
    if err := c.BodyParser(&userResponse); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid response"})
    }

    existingResponse, _ := h.responseRepo.GetByUserAndQuestionnaire(userResponse.UserID, uint(questionnaireID))
    if existingResponse != nil {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User has already voted"})
    }

    questions, _ := h.questionRepo.GetByQuestionnaireID(uint(questionnaireID))
    if len(questions) == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No questions found for this questionnaire"})
    }

    for _, answer := range userResponse.Answers {
        // Validating answer for each question (e.g., ensuring it matches the expected type)
        if err := validateAnswer(answer, questions); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
        }
    }

    err = h.responseRepo.Create(&userResponse)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to submit vote"})
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Vote submitted successfully"})
}

func validateAnswer(answer models.Answer, questions []models.Question) error {
    for _, question := range questions {
        if question.ID == answer.QuestionID {
            switch question.QuestionType {
            case "multiple_choice", "checkbox":
                if !contains(question.Options, answer.AnswerText) {
                    return fmt.Errorf("Invalid answer option for question %d", question.ID)
                }
            case "text":
                // No validation needed for text answers
            default:
                return fmt.Errorf("Unknown question type for question %d", question.ID)
            }
        }
    }
    return nil
}

func contains(options []string, answer string) bool {
    for _, option := range options {
        if option == answer {
            return true
        }
    }
    return false
}
