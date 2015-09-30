package models

import (
	"testing"
	"time"
)

func TestNewApiColumnType(t *testing.T) {
	var id int = 1
	var name = "column"
	var position int64 = 1
	var description = "description"
	var required = true
	var regexp = "^[0-9]*$"
	var horAlignmentHead = "left"
	var horAlignmentBody = "right"
	var аpiColumnType *ApiColumnType

	аpiColumnType = NewApiColumnType(id, name, position, description, required, regexp, horAlignmentHead, horAlignmentBody)
	if аpiColumnType.ID != id {
		t.Error("ID field is not properly initialized")
	}
	if аpiColumnType.Name != name {
		t.Error("Name field is not properly initialized")
	}
	if аpiColumnType.Position != position {
		t.Error("Position field is not properly initialized")
	}
	if аpiColumnType.Description != description {
		t.Error("Description field is not properly initialized")
	}
	if аpiColumnType.Required != required {
		t.Error("Required field is not properly initialized")
	}
	if аpiColumnType.Regexp != regexp {
		t.Error("Regexp field is not properly initialized")
	}
	if аpiColumnType.HorAlignmentHead != horAlignmentHead {
		t.Error("HorAlignmentHead field is not properly initialized")
	}
	if аpiColumnType.HorAlignmentBody != horAlignmentBody {
		t.Error("HorAlignmentBody field is not properly initialized")
	}
}

func TestNewDtoColumnType(t *testing.T) {
	var id int = 1
	var name = "column"
	var position int64 = 1
	var description = "description"
	var required = true
	var regexp = "^[0-9]*$"
	var horAlignmentHead Alignment = ALIGNMENT_LEFT
	var horAlignmentBody Alignment = ALIGNMNET_RIGHT
	var created = time.Now()
	var active = true
	var public = true
	var dtoColumnType *DtoColumnType

	dtoColumnType = NewDtoColumnType(id, name, position, description, required, regexp, horAlignmentHead, horAlignmentBody, created, active, public)
	if dtoColumnType.ID != id {
		t.Error("ID field is not properly initialized")
	}
	if dtoColumnType.Name != name {
		t.Error("Name field is not properly initialized")
	}
	if dtoColumnType.Position != position {
		t.Error("Position field is not properly initialized")
	}
	if dtoColumnType.Description != description {
		t.Error("Description field is not properly initialized")
	}
	if dtoColumnType.Required != required {
		t.Error("Required field is not properly initialized")
	}
	if dtoColumnType.Regexp != regexp {
		t.Error("Regexp field is not properly initialized")
	}
	if dtoColumnType.HorAlignmentHead != horAlignmentHead {
		t.Error("HorAlignmentHead field is not properly initialized")
	}
	if dtoColumnType.HorAlignmentBody != horAlignmentBody {
		t.Error("HorAlignmentBody field is not properly initialized")
	}
	if dtoColumnType.Created != created {
		t.Error("Created field is not properly initialized")
	}
	if dtoColumnType.Active != active {
		t.Error("Active field is not properly initialized")
	}
	if dtoColumnType.Public != public {
		t.Error("Public field is not properly initialized")
	}
}
