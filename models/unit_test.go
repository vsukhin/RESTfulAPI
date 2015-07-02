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
	var subscribed = true
	var paid = true
	var dtoUnit *DtoUnit

	dtoUnit = NewDtoUnit(id, created, name, active, created, created)
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
	if dtoUnit.Subscribed != subscribed {
		t.Error("Subscribed field is not properly initialized")
	}
	if dtoUnit.Paid != paid {
		t.Error("Paid field is not properly initialized")
	}
	if dtoUnit.Begin_Paid != created {
		t.Error("Paid begin field is not properly initialized")
	}
	if dtoUnit.End_Paid != created {
		t.Error("Paid end field is not properly initialized")
	}
}
