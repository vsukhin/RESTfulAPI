package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type RecognizeProductRepository interface {
	Exists(name string) (found bool, err error)
	FindByName(name string) (recognizeproduct *models.DtoRecognizeProduct, err error)
	Get(id int) (recognizeproduct *models.DtoRecognizeProduct, err error)
	GetAll() (recognizeproducts *[]models.ApiRecognizeProduct, err error)
	Create(recognizeproduct *models.DtoRecognizeProduct) (err error)
	CreateAll(recognizeproducts *[]models.DtoRecognizeProduct) (err error)
	Update(recognizeproduct *models.DtoRecognizeProduct) (err error)
	Deactivate(recognizeproduct *models.DtoRecognizeProduct) (err error)
}

type RecognizeProductService struct {
	*Repository
}

func NewRecognizeProductService(repository *Repository) *RecognizeProductService {
	repository.DbContext.AddTableWithName(models.DtoRecognizeProduct{}, repository.Table).SetKeys(true, "id")
	return &RecognizeProductService{Repository: repository}
}

func (recognizeproductservice *RecognizeProductService) Exists(name string) (found bool, err error) {
	var count int64
	count, err = recognizeproductservice.DbContext.SelectInt("select count(*) from "+recognizeproductservice.Table+" where name = ?", name)
	if err != nil {
		log.Error("Error during getting recognize product object from database %v with value %v", err, name)
		return false, err
	}

	return count != 0, nil
}

func (recognizeproductservice *RecognizeProductService) FindByName(name string) (recognizeproduct *models.DtoRecognizeProduct, err error) {
	recognizeproduct = new(models.DtoRecognizeProduct)
	err = recognizeproductservice.DbContext.SelectOne(recognizeproduct, "select * from "+recognizeproductservice.Table+" where name = ?", name)
	if err != nil {
		log.Error("Error during finding recognize product object in database %v with value %v", err, name)
		return nil, err
	}

	return recognizeproduct, nil
}

func (recognizeproductservice *RecognizeProductService) Get(id int) (recognizeproduct *models.DtoRecognizeProduct, err error) {
	recognizeproduct = new(models.DtoRecognizeProduct)
	err = recognizeproductservice.DbContext.SelectOne(recognizeproduct,
		"select * from "+recognizeproductservice.Table+" where id = ? order by position asc", id)
	if err != nil {
		log.Error("Error during getting recognize product object from database %v with value %v", err, id)
		return nil, err
	}

	return recognizeproduct, nil
}

func (recognizeproductservice *RecognizeProductService) GetAll() (recognizeproducts *[]models.ApiRecognizeProduct, err error) {
	recognizeproducts = new([]models.ApiRecognizeProduct)
	_, err = recognizeproductservice.DbContext.Select(recognizeproducts,
		"select id, position, name, description, increase from "+recognizeproductservice.Table+" where active = 1")
	if err != nil {
		log.Error("Error during getting all recognize product object from database %v", err)
		return nil, err
	}

	return recognizeproducts, nil
}

func (recognizeproductservice *RecognizeProductService) Create(recognizeproduct *models.DtoRecognizeProduct) (err error) {
	err = recognizeproductservice.DbContext.Insert(recognizeproduct)
	if err != nil {
		log.Error("Error during creating recognize product object in database %v", err)
		return err
	}

	return nil
}

func (recognizeproductservice *RecognizeProductService) CreateAll(recognizeproducts *[]models.DtoRecognizeProduct) (err error) {
	var trans *gorp.Transaction

	trans, err = recognizeproductservice.DbContext.Begin()
	if err != nil {
		log.Error("Error during creating recognize product object in database %v", err)
		return err
	}

	for i, recognizeproduct := range *recognizeproducts {
		err = trans.Insert(&recognizeproduct)
		(*recognizeproducts)[i].ID = recognizeproduct.ID
		if err != nil {
			_ = trans.Rollback()
			log.Error("Error during creating recognize product object in database %v", err)
			return err
		}
	}

	err = trans.Commit()
	if err != nil {
		log.Error("Error during creating recognize product object in database %v", err)
		return err
	}

	return nil
}

func (recognizeproductservice *RecognizeProductService) Update(recognizeproduct *models.DtoRecognizeProduct) (err error) {
	_, err = recognizeproductservice.DbContext.Update(recognizeproduct)
	if err != nil {
		log.Error("Error during updating recognize product object in database %v with value %v", err, recognizeproduct.ID)
		return err
	}

	return nil
}

func (recognizeproductservice *RecognizeProductService) Deactivate(recognizeproduct *models.DtoRecognizeProduct) (err error) {
	_, err = recognizeproductservice.DbContext.Exec("update "+recognizeproductservice.Table+" set active = 0 where id = ?", recognizeproduct.ID)
	if err != nil {
		log.Error("Error during deactivating recognize product object in database %v with value %v", err, recognizeproduct.ID)
		return err
	}

	return nil
}
