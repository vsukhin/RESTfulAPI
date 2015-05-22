package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type InputFileRepository interface {
	Get(order_id int64, file_id int64) (inputfile *models.DtoInputFile, err error)
	GetByOrder(order_id int64) (inputfiles *[]models.ApiInputFile, err error)
	Create(dtoinputfile *models.DtoInputFile, trans *gorp.Transaction) (err error)
	DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error)
}

type InputFileService struct {
	*Repository
}

func NewInputFileService(repository *Repository) *InputFileService {
	repository.DbContext.AddTableWithName(models.DtoInputFile{}, repository.Table).SetKeys(false, "order_id", "file_id")
	return &InputFileService{Repository: repository}
}

func (inputfileservice *InputFileService) Get(order_id int64, file_id int64) (inputfile *models.DtoInputFile, err error) {
	inputfile = new(models.DtoInputFile)
	err = inputfileservice.DbContext.SelectOne(inputfile, "select * from "+inputfileservice.Table+
		" where order_id = ? and file_id = ?", order_id, file_id)
	if err != nil {
		log.Error("Error during getting input file object from database %v with value %v, %v", err, order_id, file_id)
		return nil, err
	}

	return inputfile, nil
}

func (inputfileservice *InputFileService) GetByOrder(order_id int64) (inputfiles *[]models.ApiInputFile, err error) {
	inputfiles = new([]models.ApiInputFile)
	_, err = inputfileservice.DbContext.Select(inputfiles,
		"select i.file_id, f.name from "+inputfileservice.Table+" i inner join files f on i.file_id = f.id where order_id = ?", order_id)
	if err != nil {
		log.Error("Error during getting all input file object from database %v with value %v", err, order_id)
		return nil, err
	}

	return inputfiles, nil
}

func (inputfileservice *InputFileService) Create(dtoinputfile *models.DtoInputFile, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtoinputfile)
	} else {
		err = inputfileservice.DbContext.Insert(dtoinputfile)
	}
	if err != nil {
		log.Error("Error during creating input file object in database %v", err)
		return err
	}

	return nil
}

func (inputfileservice *InputFileService) DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+inputfileservice.Table+" where order_id = ?", order_id)
	} else {
		_, err = inputfileservice.DbContext.Exec("delete from "+inputfileservice.Table+" where order_id = ?", order_id)
	}
	if err != nil {
		log.Error("Error during deleting input file objects for order object in database %v with value %v", err, order_id)
		return err
	}

	return nil
}
