package administration

import (
	"application/config"
	"application/helpers"
	"application/models"
	"application/services"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
	"strings"
	"time"
	"types"
)

// get /api/v1.0/administration/classification/contacts/
func GetClassifiers(w http.ResponseWriter, request *http.Request, r render.Render,
	classifierrepository services.ClassifierRepository, session *models.DtoSession) {
	query := ""
	var filters *[]models.FilterExp
	filters, err := helpers.GetFilterArray(new(models.ClassifierSearch), nil, request, r, session.Language)
	if err != nil {
		return
	}
	if len(*filters) != 0 {
		var masks []string
		for _, filter := range *filters {
			var exps []string
			for _, field := range filter.Fields {
				exps = append(exps, field+" "+filter.Op+" "+filter.Value)
			}
			masks = append(masks, "("+strings.Join(exps, " or ")+")")
		}
		query += " where "
		query += strings.Join(masks, " and ")
	}

	var sorts *[]models.OrderExp
	sorts, err = helpers.GetOrderArray(new(models.ClassifierSearch), request, r, session.Language)
	if err != nil {
		return
	}
	if len(*sorts) != 0 {
		var orders []string
		for _, sort := range *sorts {
			orders = append(orders, " "+sort.Field+" "+sort.Order)
		}
		query += " order by"
		query += strings.Join(orders, ",")
	}

	var limit string
	limit, err = helpers.GetLimitQuery(request, r, session.Language)
	if err != nil {
		return
	}
	query += limit

	classifiers, err := classifierrepository.GetAll(query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(classifiers, len(*classifiers), w, r)
}

// post /api/v1.0/administration/classification/contacts/
func CreateClassifier(errors binding.Errors, veiwclassifier models.ViewClassifier, r render.Render,
	classifierrepository services.ClassifierRepository, session *models.DtoSession) {
	if helpers.CheckValidation(&veiwclassifier, errors, r, session.Language) != nil {
		return
	}

	dtoclassifier := new(models.DtoClassifier)
	dtoclassifier.Name = veiwclassifier.Name
	dtoclassifier.Created = time.Now()
	dtoclassifier.Active = true

	err := classifierrepository.Create(dtoclassifier)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongClassifier(dtoclassifier.ID, dtoclassifier.Name, !dtoclassifier.Active))
}

// get /api/v1.0/administration/classification/contacts/:id/
func GetClassifier(r render.Render, params martini.Params, classifierrepository services.ClassifierRepository, session *models.DtoSession) {
	classifierid, err := helpers.CheckParameterInt(r, params[helpers.PARAM_NAME_CLASSIFIER_ID], session.Language)
	if err != nil {
		return
	}

	dtoclassifier, err := classifierrepository.Get(int(classifierid))
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongClassifier(dtoclassifier.ID, dtoclassifier.Name, !dtoclassifier.Active))
}

// put /api/v1.0/administration/classification/contacts/:id/
func UpdateClassifier(errors binding.Errors, veiwclassifier models.ViewUpdateClassifier, r render.Render, params martini.Params,
	classifierrepository services.ClassifierRepository, session *models.DtoSession) {
	if helpers.CheckValidation(&veiwclassifier, errors, r, session.Language) != nil {
		return
	}
	classifierid, err := helpers.CheckParameterInt(r, params[helpers.PARAM_NAME_CLASSIFIER_ID], session.Language)
	if err != nil {
		return
	}

	dtoclassifier, err := classifierrepository.Get(int(classifierid))
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	dtoclassifier.Name = veiwclassifier.Name
	dtoclassifier.Active = !veiwclassifier.Deleted
	err = classifierrepository.Update(dtoclassifier)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongClassifier(dtoclassifier.ID, dtoclassifier.Name, !dtoclassifier.Active))
}

// delete /api/v1.0/administration/classification/contacts/:id/
func DeleteClassifier(r render.Render, params martini.Params, classifierrepository services.ClassifierRepository, session *models.DtoSession) {
	classifierid, err := helpers.CheckParameterInt(r, params[helpers.PARAM_NAME_CLASSIFIER_ID], session.Language)
	if err != nil {
		return
	}

	classifier, err := classifierrepository.Get(int(classifierid))
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	if !classifier.Active {
		log.Error("Classifer is not active %v", classifier.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	err = classifierrepository.Deactivate(classifier)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}
