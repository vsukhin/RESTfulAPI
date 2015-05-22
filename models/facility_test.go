package models

import (
	"testing"
	"time"
)

func TestApiShortFacility(t *testing.T) {
	var id int64 = 1
	var name = "facility"
	var description = "description"
	var аpiFacility *ApiShortFacility

	аpiFacility = NewApiShortFacility(id, name, description)
	if аpiFacility.ID != id {
		t.Error("ID field is not properly initialized")
	}
	if аpiFacility.Name != name {
		t.Error("Name field is not properly initialized")
	}
	if аpiFacility.Description != description {
		t.Error("Description field is not properly initialized")
	}
}

func TestApiLongFacility(t *testing.T) {
	var id int64 = 1
	var name = "facility"
	var description = "description"
	var active = true
	var аpiFacility *ApiLongFacility

	аpiFacility = NewApiLongFacility(id, name, description, active)
	if аpiFacility.ID != id {
		t.Error("ID field is not properly initialized")
	}
	if аpiFacility.Name != name {
		t.Error("Name field is not properly initialized")
	}
	if аpiFacility.Description != description {
		t.Error("Description field is not properly initialized")
	}
	if аpiFacility.Active != active {
		t.Error("Active field is not properly initialized")
	}
}

func TestNewApiFullFacility(t *testing.T) {
	var id int64 = 1
	var category_id = 1
	var alias = "alias"
	var name = "facility"
	var description = "description"
	var descriptionsoon = "description soon"
	var active = true
	var picNormal_id int64 = 1
	var picOver_id int64 = 1
	var picSoon_id int64 = 1
	var picDisable_id int64 = 1
	var apiFullFacility *ApiFullFacility

	apiFullFacility = NewApiFullFacility(id, category_id, alias, name, description, descriptionsoon, active,
		picNormal_id, picOver_id, picSoon_id, picDisable_id)
	if apiFullFacility.ID != id {
		t.Error("ID field is not properly initialized")
	}
	if apiFullFacility.Category_ID != category_id {
		t.Error("Category field is not properly initialized")
	}
	if apiFullFacility.Alias != alias {
		t.Error("Alias field is not properly initialized")
	}
	if apiFullFacility.Name != name {
		t.Error("Name field is not properly initialized")
	}
	if apiFullFacility.Description != description {
		t.Error("Description field is not properly initialized")
	}
	if apiFullFacility.DescriptionSoon != descriptionsoon {
		t.Error("Description soon field is not properly initialized")
	}
	if apiFullFacility.Active != active {
		t.Error("Active field is not properly initialized")
	}
	if apiFullFacility.PicNormal_ID != picNormal_id {
		t.Error("Picture normal id field is not properly initialized")
	}
	if apiFullFacility.PicOver_ID != picOver_id {
		t.Error("Picture over id field is not properly initialized")
	}
	if apiFullFacility.PicSoon_ID != picSoon_id {
		t.Error("Picture soon id field is not properly initialized")
	}
	if apiFullFacility.PicDisable_ID != picDisable_id {
		t.Error("Picture disable id field is not properly initialized")
	}
}

func TestNewDtoFacility(t *testing.T) {
	var id int64 = 1
	var category_id = 1
	var name = "facility"
	var description = "description"
	var descriptionsoon = "description soon"
	var created = time.Now()
	var active = true
	var picNormal_id int64 = 1
	var picOver_id int64 = 1
	var picSoon_id int64 = 1
	var picDisable_id int64 = 1
	var alias = "alias"
	var dtoFacility *DtoFacility

	dtoFacility = NewDtoFacility(id, name, description, created, active, category_id, descriptionsoon,
		picNormal_id, picOver_id, picSoon_id, picDisable_id, alias)
	if dtoFacility.ID != id {
		t.Error("ID field is not properly initialized")
	}
	if dtoFacility.Category_ID != category_id {
		t.Error("Category ID field is not properly initialized")
	}
	if dtoFacility.Name != name {
		t.Error("Name field is not properly initialized")
	}
	if dtoFacility.Description != description {
		t.Error("Description field is not properly initialized")
	}
	if dtoFacility.DescriptionSoon != descriptionsoon {
		t.Error("Description soon field is not properly initialized")
	}
	if dtoFacility.Created != created {
		t.Error("Created field is not properly initialized")
	}
	if dtoFacility.Active != active {
		t.Error("Active field is not properly initialized")
	}
	if dtoFacility.PicNormal_ID != picNormal_id {
		t.Error("Picture normal id field is not properly initialized")
	}
	if dtoFacility.PicOver_ID != picOver_id {
		t.Error("Picture over id field is not properly initialized")
	}
	if dtoFacility.PicSoon_ID != picSoon_id {
		t.Error("Picture soon id field is not properly initialized")
	}
	if dtoFacility.PicDisable_ID != picDisable_id {
		t.Error("Picture disable id field is not properly initialized")
	}
	if dtoFacility.Alias != alias {
		t.Error("Alias field is not properly initialized")
	}
}
