package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"errors"
	"github.com/martini-contrib/render"
	"net/http"
	"types"
)

func CheckVerifyProductInternal(verifyproduct_id int, verifyproductrepository services.VerifyProductRepository) (
	dtoverifyproduct *models.DtoVerifyProduct, err error) {
	dtoverifyproduct, err = verifyproductrepository.Get(verifyproduct_id)
	if err != nil {
		return nil, err
	}
	if !dtoverifyproduct.Active {
		log.Error("Verify product is not active %v", dtoverifyproduct.ID)
		return nil, errors.New("Not active verify product")
	}

	return dtoverifyproduct, nil
}

func CheckVerifyProduct(verifyproduct_id int, r render.Render, verifyproductrepository services.VerifyProductRepository,
	language string) (dtoverifyproduct *models.DtoVerifyProduct, err error) {
	dtoverifyproduct, err = CheckVerifyProductInternal(verifyproduct_id, verifyproductrepository)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	return dtoverifyproduct, nil
}
