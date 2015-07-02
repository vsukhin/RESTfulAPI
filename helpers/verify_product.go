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

func CheckVerifyProduct(verifyproduct_id int, r render.Render, verifyproductrepository services.VerifyProductRepository,
	language string) (dtoverifyproduct *models.DtoVerifyProduct, err error) {
	dtoverifyproduct, err = verifyproductrepository.Get(verifyproduct_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if !dtoverifyproduct.Active {
		log.Error("Verify product is not active %v", dtoverifyproduct.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Not active verify product")
	}

	return dtoverifyproduct, nil
}
