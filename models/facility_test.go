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

func TestNewDtoFacility(t *testing.T) {
	var id int64 = 1
	var name = "facility"
	var description = "description"
	var created = time.Now()
	var active = true
	var dtoFacility *DtoFacility

	dtoFacility = NewDtoFacility(id, name, description, created, active)
	if dtoFacility.ID != id {
		t.Error("ID field is not properly initialized")
	}
	if dtoFacility.Name != name {
		t.Error("Name field is not properly initialized")
	}
	if dtoFacility.Description != description {
		t.Error("Description field is not properly initialized")
	}
	if dtoFacility.Created != created {
		t.Error("Created field is not properly initialized")
	}
	if dtoFacility.Active != active {
		t.Error("Active field is not properly initialized")
	}
}
