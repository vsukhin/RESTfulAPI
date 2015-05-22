package models

import (
	"testing"
)

func TestNewApiGroup(t *testing.T) {
	var id int = 1
	var name = "group"
	var аpiGroup *ApiGroup

	аpiGroup = NewApiGroup(id, name)
	if аpiGroup.ID != id {
		t.Error("ID field is not properly initialized")
	}
	if аpiGroup.Name != name {
		t.Error("Name field is not properly initialized")
	}
}

func TestNewDtoGroup(t *testing.T) {
	var id int = 1
	var name = "group"
	var isdefault = true
	var isactive = true
	var dtoGroup *DtoGroup

	dtoGroup = NewDtoGroup(id, name, isdefault, isactive)
	if dtoGroup.ID != id {
		t.Error("ID field is not properly initialized")
	}
	if dtoGroup.Name != name {
		t.Error("Name field is not properly initialized")
	}
	if dtoGroup.IsDefault != isdefault {
		t.Error("IsDefault field is not properly initialized")
	}
	if dtoGroup.IsActive != isactive {
		t.Error("IsActive field is not properly initialized")
	}
}
