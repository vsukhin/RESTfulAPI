package services

import (
	"application/models"
)

type EventRepository interface {
	Get(id int) (event *models.DtoEvent, err error)
	GetAll() (events *[]models.ApiEvent, err error)
}

type EventService struct {
	*Repository
}

func NewEventService(repository *Repository) *EventService {
	repository.DbContext.AddTableWithName(models.DtoEvent{}, repository.Table).SetKeys(true, "id")
	return &EventService{Repository: repository}
}

func (eventservice *EventService) Get(id int) (event *models.DtoEvent, err error) {
	event = new(models.DtoEvent)
	err = eventservice.DbContext.SelectOne(event, "select * from "+eventservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting event object from database %v with value %v", err, id)
		return nil, err
	}

	return event, nil
}

func (eventservice *EventService) GetAll() (events *[]models.ApiEvent, err error) {
	events = new([]models.ApiEvent)
	_, err = eventservice.DbContext.Select(events, "select id, name, position from "+eventservice.Table+" where active = 1 order by position asc")
	if err != nil {
		log.Error("Error during getting all event object from database %v", err)
		return nil, err
	}

	return events, nil
}
