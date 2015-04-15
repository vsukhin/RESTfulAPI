package services

import (
	"application/models"
)

type ClassifierRepository interface {
	Get(id int) (classifier *models.DtoClassifier, err error)
	GetAllAvailable() (classifiers *[]models.ApiShortClassifier, err error)
	GetAll(filter string) (classifiers *[]models.ApiLongClassifier, err error)
	Create(classifier *models.DtoClassifier) (err error)
	Update(classifier *models.DtoClassifier) (err error)
	Deactivate(classifier *models.DtoClassifier) (err error)
}

type ClassifierService struct {
	*Repository
}

func NewClassifierService(repository *Repository) *ClassifierService {
	repository.DbContext.AddTableWithName(models.DtoClassifier{}, repository.Table).SetKeys(true, "id")
	return &ClassifierService{Repository: repository}
}

func (classifierservice *ClassifierService) Get(id int) (classifier *models.DtoClassifier, err error) {
	classifier = new(models.DtoClassifier)
	err = classifierservice.DbContext.SelectOne(classifier, "select * from "+classifierservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting classifier object from database %v with value %v", err, id)
		return nil, err
	}

	return classifier, nil
}

func (classifierservice *ClassifierService) GetAllAvailable() (classifiers *[]models.ApiShortClassifier, err error) {
	classifiers = new([]models.ApiShortClassifier)
	_, err = classifierservice.DbContext.Select(classifiers, "select id, name from "+classifierservice.Table+" where active = 1")
	if err != nil {
		log.Error("Error during getting all classifier object from database %v", err)
		return nil, err
	}

	return classifiers, nil
}

func (classifierservice *ClassifierService) GetAll(filter string) (classifiers *[]models.ApiLongClassifier, err error) {
	classifiers = new([]models.ApiLongClassifier)
	_, err = classifierservice.DbContext.Select(classifiers, "select id, name, not active as del from "+classifierservice.Table+filter)
	if err != nil {
		log.Error("Error during getting all classifier object from database %v", err)
		return nil, err
	}

	return classifiers, nil
}

func (classifierservice *ClassifierService) Create(classifier *models.DtoClassifier) (err error) {
	err = classifierservice.DbContext.Insert(classifier)
	if err != nil {
		log.Error("Error during creating classifier object in database %v", err)
		return err
	}

	return nil
}

func (classifierservice *ClassifierService) Update(classifier *models.DtoClassifier) (err error) {
	_, err = classifierservice.DbContext.Update(classifier)
	if err != nil {
		log.Error("Error during updating classifier object in database %v with value %v", err, classifier.ID)
		return err
	}

	return nil
}

func (classifierservice *ClassifierService) Deactivate(classifier *models.DtoClassifier) (err error) {
	_, err = classifierservice.DbContext.Exec("update "+classifierservice.Table+" set active = 0 where id = ?", classifier.ID)
	if err != nil {
		log.Error("Error during deactivating classifier object in database %v with value %v", err, classifier.ID)
		return err
	}

	return nil
}
