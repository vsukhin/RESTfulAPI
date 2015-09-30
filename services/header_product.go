package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type HeaderProductRepository interface {
	Exists(name string) (found bool, err error)
	FindByName(name string) (headerproduct *models.DtoHeaderProduct, err error)
	Get(id int) (headerproduct *models.DtoHeaderProduct, err error)
	GetAll() (headerproducts *[]models.ApiHeaderProduct, err error)
	Create(headerproduct *models.DtoHeaderProduct) (err error)
	CreateAll(headerproducts *[]models.DtoHeaderProduct) (err error)
	Update(headerproduct *models.DtoHeaderProduct) (err error)
	Deactivate(headerproduct *models.DtoHeaderProduct) (err error)
}

type HeaderProductService struct {
	*Repository
}

func NewHeaderProductService(repository *Repository) *HeaderProductService {
	repository.DbContext.AddTableWithName(models.DtoHeaderProduct{}, repository.Table).SetKeys(true, "id")
	return &HeaderProductService{Repository: repository}
}

func (headerproductservice *HeaderProductService) Exists(name string) (found bool, err error) {
	var count int64
	count, err = headerproductservice.DbContext.SelectInt("select count(*) from "+headerproductservice.Table+" where name = ?", name)
	if err != nil {
		log.Error("Error during getting header product object from database %v with value %v", err, name)
		return false, err
	}

	return count != 0, nil
}

func (headerproductservice *HeaderProductService) FindByName(name string) (headerproduct *models.DtoHeaderProduct, err error) {
	headerproduct = new(models.DtoHeaderProduct)
	err = headerproductservice.DbContext.SelectOne(headerproduct, "select * from "+headerproductservice.Table+" where name = ?", name)
	if err != nil {
		log.Error("Error during finding header product object in database %v with value %v", err, name)
		return nil, err
	}

	return headerproduct, nil
}

func (headerproductservice *HeaderProductService) Get(id int) (headerproduct *models.DtoHeaderProduct, err error) {
	headerproduct = new(models.DtoHeaderProduct)
	err = headerproductservice.DbContext.SelectOne(headerproduct,
		"select * from "+headerproductservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting header product object from database %v with value %v", err, id)
		return nil, err
	}

	return headerproduct, nil
}

func (headerproductservice *HeaderProductService) GetAll() (headerproducts *[]models.ApiHeaderProduct, err error) {
	headerproducts = new([]models.ApiHeaderProduct)
	_, err = headerproductservice.DbContext.Select(headerproducts,
		"select id, position, name, description, increase, fee_once, fee_monthly from "+headerproductservice.Table+
			" where active = 1 order by position asc")
	if err != nil {
		log.Error("Error during getting all header product object from database %v", err)
		return nil, err
	}

	return headerproducts, nil
}

func (headerproductservice *HeaderProductService) Create(headerproduct *models.DtoHeaderProduct) (err error) {
	err = headerproductservice.DbContext.Insert(headerproduct)
	if err != nil {
		log.Error("Error during creating header product object in database %v", err)
		return err
	}

	return nil
}

func (headerproductservice *HeaderProductService) CreateAll(headerproducts *[]models.DtoHeaderProduct) (err error) {
	var trans *gorp.Transaction

	trans, err = headerproductservice.DbContext.Begin()
	if err != nil {
		log.Error("Error during creating header product object in database %v", err)
		return err
	}

	for i, headerproduct := range *headerproducts {
		err = trans.Insert(&headerproduct)
		(*headerproducts)[i].ID = headerproduct.ID
		if err != nil {
			_ = trans.Rollback()
			log.Error("Error during creating header product object in database %v", err)
			return err
		}
	}

	err = trans.Commit()
	if err != nil {
		log.Error("Error during creating header product object in database %v", err)
		return err
	}

	return nil
}

func (headerproductservice *HeaderProductService) Update(headerproduct *models.DtoHeaderProduct) (err error) {
	_, err = headerproductservice.DbContext.Update(headerproduct)
	if err != nil {
		log.Error("Error during updating header product object in database %v with value %v", err, headerproduct.ID)
		return err
	}

	return nil
}

func (headerproductservice *HeaderProductService) Deactivate(headerproduct *models.DtoHeaderProduct) (err error) {
	_, err = headerproductservice.DbContext.Exec("update "+headerproductservice.Table+" set active = 0 where id = ?", headerproduct.ID)
	if err != nil {
		log.Error("Error during deactivating header product object in database %v with value %v", err, headerproduct.ID)
		return err
	}

	return nil
}
