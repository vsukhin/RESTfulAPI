package models

import (
	"testing"
)

func TestApiDataFormat(t *testing.T) {
	var id int = 1
	var name = "format"
	var description = "description"
	var аpiDataFormat *ApiDataFormat

	аpiDataFormat = NewApiDataFormat(id, name, description)
	if аpiDataFormat.ID != id {
		t.Error("ID field is not properly initialized")
	}
	if аpiDataFormat.Name != name {
		t.Error("Name field is not properly initialized")
	}
	if аpiDataFormat.Description != description {
		t.Error("Description field is not properly initialized")
	}
}

func TestNewDtoDataFormat(t *testing.T) {
	var id int = 1
	var name = "format"
	var description = "description"
	var dtoDataFormat *DtoDataFormat

	dtoDataFormat = NewDtoDataFormat(id, name, description)
	if dtoDataFormat.ID != id {
		t.Error("ID field is not properly initialized")
	}
	if dtoDataFormat.Name != name {
		t.Error("Name field is not properly initialized")
	}
	if dtoDataFormat.Description != description {
		t.Error("Description field is not properly initialized")
	}
}
