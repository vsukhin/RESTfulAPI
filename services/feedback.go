package services

import (
	"application/models"
)

type FeedbackRepository interface {
	Create(dtofeedback *models.DtoFeedback) (err error)
}

type FeedbackService struct {
	*Repository
}

func NewFeedbackService(repository *Repository) *FeedbackService {
	repository.DbContext.AddTableWithName(models.DtoFeedback{}, repository.Table).SetKeys(true, "id")
	return &FeedbackService{
		repository,
	}
}

func (feedbackservice *FeedbackService) Create(dtofeedback *models.DtoFeedback) (err error) {
	err = feedbackservice.DbContext.Insert(dtofeedback)
	if err != nil {
		log.Error("Error during creating feedback object in database %v", err)
		return err
	}

	return nil
}
