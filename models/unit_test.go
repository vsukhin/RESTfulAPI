package models

import (
	"testing"
	"time"
)

func TestNewDtoUnit(t *testing.T) {
	var id int64 = 1
	var created = time.Now()
	var name = "Name"
	var active = true
	var dtoUnit *DtoUnit

	dtoUnit = NewDtoUnit(id, created, name, active)
	if dtoUnit.ID != id {
		t.Error("ID field is not properly initialized")
	}
	if dtoUnit.Created != created {
		t.Error("Created field is not properly initialized")
	}
	if dtoUnit.Name != name {
		t.Error("Name field is not properly initialized")
	}
	if dtoUnit.Active != active {
		t.Error("Active field is not properly initialized")
	}
}
