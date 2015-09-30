package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type VerifyProductRepository interface {
	Exists(name string) (found bool, err error)
	FindByName(name string) (verifyproduct *models.DtoVerifyProduct, err error)
	Get(id int) (verifyproduct *models.DtoVerifyProduct, err error)
	GetAll() (verifyproducts *[]models.ApiVerifyProduct, err error)
	Create(verifyproduct *models.DtoVerifyProduct) (err error)
	CreateAll(verifyproducts *[]models.DtoVerifyProduct) (err error)
	Update(verifyproduct *models.DtoVerifyProduct) (err error)
	Deactivate(verifyproduct *models.DtoVerifyProduct) (err error)
}

type VerifyProductService struct {
	*Repository
}

func NewVerifyProductService(repository *Repository) *VerifyProductService {
	repository.DbContext.AddTableWithName(models.DtoVerifyProduct{}, repository.Table).SetKeys(true, "id")
	return &VerifyProductService{Repository: repository}
}

func (verifyproductservice *VerifyProductService) Exists(name string) (found bool, err error) {
	var count int64
	count, err = verifyproductservice.DbContext.SelectInt("select count(*) from "+verifyproductservice.Table+" where name = ?", name)
	if err != nil {
		log.Error("Error during getting verify product object from database %v with value %v", err, name)
		return false, err
	}

	return count != 0, nil
}

func (verifyproductservice *VerifyProductService) FindByName(name string) (verifyproduct *models.DtoVerifyProduct, err error) {
	verifyproduct = new(models.DtoVerifyProduct)
	err = verifyproductservice.DbContext.SelectOne(verifyproduct, "select * from "+verifyproductservice.Table+" where name = ?", name)
	if err != nil {
		log.Error("Error during finding verify product object in database %v with value %v", err, name)
		return nil, err
	}

	return verifyproduct, nil
}

func (verifyproductservice *VerifyProductService) Get(id int) (verifyproduct *models.DtoVerifyProduct, err error) {
	verifyproduct = new(models.DtoVerifyProduct)
	err = verifyproductservice.DbContext.SelectOne(verifyproduct,
		"select * from "+verifyproductservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting verify product object from database %v with value %v", err, id)
		return nil, err
	}

	return verifyproduct, nil
}

func (verifyproductservice *VerifyProductService) GetAll() (verifyproducts *[]models.ApiVerifyProduct, err error) {
	verifyproducts = new([]models.ApiVerifyProduct)
	_, err = verifyproductservice.DbContext.Select(verifyproducts,
		"select id, position, name, description from "+verifyproductservice.Table+" where active = 1 order by position asc")
	if err != nil {
		log.Error("Error during getting all verify product object from database %v", err)
		return nil, err
	}

	return verifyproducts, nil
}

func (verifyproductservice *VerifyProductService) Create(verifyproduct *models.DtoVerifyProduct) (err error) {
	err = verifyproductservice.DbContext.Insert(verifyproduct)
	if err != nil {
		log.Error("Error during creating verify product object in database %v", err)
		return err
	}

	return nil
}

func (verifyproductservice *VerifyProductService) CreateAll(verifyproducts *[]models.DtoVerifyProduct) (err error) {
	var trans *gorp.Transaction

	trans, err = verifyproductservice.DbContext.Begin()
	if err != nil {
		log.Error("Error during creating verify product object in database %v", err)
		return err
	}

	for i, verifyproduct := range *verifyproducts {
		err = trans.Insert(&verifyproduct)
		(*verifyproducts)[i].ID = verifyproduct.ID
		if err != nil {
			_ = trans.Rollback()
			log.Error("Error during creating verify product object in database %v", err)
			return err
		}
	}

	err = trans.Commit()
	if err != nil {
		log.Error("Error during creating verify product object in database %v", err)
		return err
	}

	return nil
}

func (verifyproductservice *VerifyProductService) Update(verifyproduct *models.DtoVerifyProduct) (err error) {
	_, err = verifyproductservice.DbContext.Update(verifyproduct)
	if err != nil {
		log.Error("Error during updating verify product object in database %v with value %v", err, verifyproduct.ID)
		return err
	}

	return nil
}

func (verifyproductservice *VerifyProductService) Deactivate(verifyproduct *models.DtoVerifyProduct) (err error) {
	_, err = verifyproductservice.DbContext.Exec("update "+verifyproductservice.Table+" set active = 0 where id = ?", verifyproduct.ID)
	if err != nil {
		log.Error("Error during deactivating verify product object in database %v with value %v", err, verifyproduct.ID)
		return err
	}

	return nil
}
