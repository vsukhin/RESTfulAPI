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

func CheckRecognizeProductInternal(recognizeproduct_id int, recognizeproductrepository services.RecognizeProductRepository) (
	dtorecognizeproduct *models.DtoRecognizeProduct, err error) {
	dtorecognizeproduct, err = recognizeproductrepository.Get(recognizeproduct_id)
	if err != nil {
		return nil, err
	}
	if !dtorecognizeproduct.Active {
		log.Error("Recognize product is not active %v", dtorecognizeproduct.ID)
		return nil, errors.New("Not active recognize product")
	}

	return dtorecognizeproduct, nil
}

func CheckRecognizeProduct(recognizeproduct_id int, r render.Render, recognizeproductrepository services.RecognizeProductRepository,
	language string) (dtorecognizeproduct *models.DtoRecognizeProduct, err error) {
	dtorecognizeproduct, err = CheckRecognizeProductInternal(recognizeproduct_id, recognizeproductrepository)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	return dtorecognizeproduct, nil
}
