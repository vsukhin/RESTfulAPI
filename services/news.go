package services

import (
	"application/models"
)

type NewsRepository interface {
	Get(id int64) (news *models.DtoNews, err error)
	GetAll(language string, count int64) (news *[]models.DtoNews, err error)
}

type NewsService struct {
	*Repository
}

func NewNewsService(repository *Repository) *NewsService {
	repository.DbContext.AddTableWithName(models.DtoNews{}, repository.Table).SetKeys(true, "id")
	return &NewsService{
		repository,
	}
}

func (newsservice *NewsService) Get(id int64) (news *models.DtoNews, err error) {
	news = new(models.DtoNews)
	err = newsservice.DbContext.SelectOne(news, "select * from "+newsservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting news object from database %v with value %v", err, id)
		return nil, err
	}

	return news, nil
}

func (newsservice *NewsService) GetAll(language string, count int64) (news *[]models.DtoNews, err error) {
	news = new([]models.DtoNews)
	_, err = newsservice.DbContext.Select(news, "select title, description from "+newsservice.Table+
		" where active = 1 and language = ? order by created desc limit ?", language, count)
	if err != nil {
		log.Error("Error during getting all news object from database %v", err)
		return nil, err
	}

	return news, nil
}
