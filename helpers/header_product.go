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

func CheckHeaderProductInternal(headerproduct_id int, headerproductrepository services.HeaderProductRepository) (
	dtoheaderproduct *models.DtoHeaderProduct, err error) {
	dtoheaderproduct, err = headerproductrepository.Get(headerproduct_id)
	if err != nil {
		return nil, err
	}
	if !dtoheaderproduct.Active {
		log.Error("Header product is not active %v", dtoheaderproduct.ID)
		return nil, errors.New("Not active header product")
	}

	return dtoheaderproduct, nil
}

func CheckHeaderProduct(headerproduct_id int, r render.Render, headerproductrepository services.HeaderProductRepository,
	language string) (dtoheaderproduct *models.DtoHeaderProduct, err error) {
	dtoheaderproduct, err = CheckHeaderProductInternal(headerproduct_id, headerproductrepository)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	return dtoheaderproduct, nil
}
